package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
	"bytes"
	"context"
	"strings"
)

type ExternalRunner struct {
	CmdName string         //binary that runs code 
	Args    []string       //default args
	Timeout time.Duration  //command timeout
	Env []string           //nil to use os.Environ()
	WorkDir string         //working dir
}

func (r *ExternalRunner) Run(code string) (string, string, error) {
	ctx := context.Background() 
	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	//create temporary file
	tmp, err := os.CreateTemp("", "snippet*")
	if err != nil {
		errRun("create temporary file", err)
	}

	tmpName := tmp.Name()

	//cleanup then close once fn ends
	defer func() {
		tmp.Close()
		_ = os.Remove(tmpName)
	}()

	switch (r.CmdName) {
	case "java":
		fileName := strings.ReplaceAll(tmpName, "/tmp/", "")
//		class := strings.Split(fileName, ".")[0]
		class := fileName
		code = fmt.Sprintf("public class %s {\n%s\n}\n", class, code)
		fmt.Println("appending Java header")
	default:
		fmt.Printf("")
	}

	fmt.Printf(code)

	//write code into temporary file
	_, err = io.WriteString(tmp, code)
	if err != nil {
		return errRun("write to temporary file", err)
	}

	//make sure the file is available
	err = tmp.Sync()
	if err != nil {
		return errRun("sync temporary file", err)
	}

	//close the file
	err = tmp.Close()
	if err != nil {
		return errRun("close temporary file", err)
	}

	//add provided args
	args := append([]string{}, r.Args...)
	//add file name to args
	args = append(args, tmpName)
	//create command
	cmd := exec.CommandContext(ctx, r.CmdName, args...)
	
	//use working dir provided if exists
	if r.WorkDir != "" {
		cmd.Dir = r.WorkDir
	}

	//use environment provided if exists
	if r.Env != nil {
		cmd.Env = r.Env
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	stderrStr := stderr.String()
	stdoutStr := stdout.String()
	
	//context exceeds timeout, classify as such
	if ctx.Err() == context.DeadlineExceeded {
		errRet := fmt.Errorf("exec exceeded timeout")
		return stdoutStr, stderrStr, errRet 
	}

	if err != nil {
		errRet := fmt.Errorf("exec:  %w", err)
		return stdoutStr, stderrStr, errRet 
	}

	
	return stdoutStr, stderrStr, nil 
}

func parseAndRun(src string, registry map[string]Runner) (string, error) {
	var out strings.Builder
	i := 0
	n := len(src)
	for {
		rel := strings.Index(src[i:], "<$")
		if rel < 0 {
			out.WriteString(src[i:])
			break
		}
		start := i + rel
		out.WriteString(src[i:start])

		// If "<?" is at very end, treat literally
		if start+2 >= n {
			out.WriteString("<$")
			i = start + 2
			continue
		}

		j := start + 2
		k := j
		for k < n {
			c := src[k]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				k++
				continue
			}
			break
		}
		if k == j {
			out.WriteString("<$")
			i = j
			continue
		}
		lang := src[j:k]

		codeStart := k
		if codeStart < n && src[codeStart] == ' ' {
			codeStart++
		}
		if codeStart >= n {
			out.WriteString(src[start:])
			break
		}

		endRel := strings.Index(src[codeStart:], "$>")
		if endRel < 0 {
			out.WriteString(src[start:])
			break
		}
		end := codeStart + endRel
		code := src[codeStart:end]

		r, ok := registry[strings.ToLower(lang)]
		if !ok {
			return "", fmt.Errorf("unknown language tag: %s", lang)
		}

		stdout, stderr, err := r.Run(code)
		if err != nil {
			return "", fmt.Errorf("runner %s failed: %w; stderr=%s", lang, err, stderr)
		}
		out.WriteString(stdout)
		i = end + 2
	}
	return out.String(), nil
}

func errRun(str string, err error) (string, string, error) {
	return "", "", fmt.Errorf("%s:  %w", str, err)
}

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
	"path/filepath"
)

type ExternalRunner struct {
	CmdName string         //binary that runs code 
	Args    []string       //default args
	Timeout time.Duration  //command timeout
	Env []string           //nil to use os.Environ()
	WorkDir string         //working dir
}

func (r *ExternalRunner) Run(code string, tmp *os.File) (string, string, error) {
	var err error 
	ctx := context.Background() 
	if r.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.Timeout)
		defer cancel()
	}

	tmpName := tmp.Name()

	//cleanup then close once fn ends
	defer func() {
		tmp.Close()
		_ = os.Remove(tmpName)
	}()
	
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

		//create temporary dir
		tmpDir, err := os.MkdirTemp("", "snippet*")
		if err != nil {
			errRun("create temporary file", err)
		}

		defer os.RemoveAll(tmpDir)

		tmp := prepForLangsWithOddReqs(lang, tmpDir)

		code = formatCode(code, lang, tmp.Name(), tmpDir)
		
		stdout, stderr, err := r.Run(code, tmp)
		if err != nil {
			return "", fmt.Errorf("runner %s failed: %w; stderr=%s", lang, err, stderr)
		}
		out.WriteString(stdout)
		i = end + 2
	}
	return out.String(), nil
}

func getCodeBetween(code string, start string, end string) string {
	star := strings.Index(code, start)
	if star == -1 {
		return code
	}
	star += len (start)

	en := strings.Index(code[star:], end)
	if en == -1 {
		return code
	}
	
	return code[star : star+en]
}

func stripImps(code string) string {
	start := strings.Index(code, "??>")
	if start == -1 {
		return code
	}
	start += 3
	return code[start:]
}

func prepForLangsWithOddReqs(lang string, tmpDir string) *os.File {
	switch (lang) {
	case "go":
		modCont := []byte("module elh\n\ngo 1.20\n")
		err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), modCont, 0644)
		if err != nil {
			fmt.Errorf("err prepping for Go:  %v\n", err)
			return nil
		}
		file, err := os.Create(filepath.Join(tmpDir, "main.go"))
		if err != nil {
			fmt.Errorf("%v", err)
			return nil
		}
		return file
	default:
		fileName := fmt.Sprintf(filepath.Base(tmpDir))
		file, err := os.Create(filepath.Join(tmpDir, fileName))
		if err != nil {
			fmt.Errorf("%v", err)
			return nil
		}
		return file
	}
}

func formatCode(code string, lang string, tmpName string, tmpDir string) string { 
	imps := getCodeBetween(code, "<??imps", "??>")
	code = stripImps(code)
	switch (lang) {
	case "java":
		fileName := strings.ReplaceAll(tmpName, fmt.Sprintf("%s/", tmpDir), "")
		class := fileName
		code = fmt.Sprintf("public class %s {\n%s\n}\n", class, code)
	case "go":
		head := fmt.Sprintf("package main\nimport (\n%s\n)\n", imps)
		code = fmt.Sprintf("%s%s", head, code)
	case "php":
		code = fmt.Sprintf("<?php\n%s\n?>", code)
	default:
	}
	return code
}

func errRun(str string, err error) (string, string, error) {
	return "", "", fmt.Errorf("%s:  %w", str, err)
}

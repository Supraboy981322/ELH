package ELH

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
	"bytes"
	"context"
	"strings"
	"net/http"
	"path/filepath"
)

type ExternalRunner struct {
	CmdName string         //binary that runs code 
	Args    []string       //default args
	Timeout time.Duration  //command timeout
	Env []string           //nil to use os.Environ()
	WorkDir string         //working dir
	Func func(runOpts RunOpts) (stdout, stderr string, err error)
}

func (r *ExternalRunner) GetRunner() *ExternalRunner {
	return r
}

func (r *ExternalRunner) Run(runOpts RunOpts) (string, string, error) {
	code := runOpts.Code
	tmp := runOpts.TmpFile

	if r.Func != nil {
		stdout, stderr, err := r.Func(runOpts)
		return stdout, stderr, err
	}
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

	var stdout_buf, stderr_buf bytes.Buffer
	var stderrWr, stdoutWr io.Writer
	stdoutWr = &stdout_buf
	stderrWr = &stderr_buf
	if Server.RenderStderr {
		stdoutWr = io.MultiWriter(stderrWr, stdoutWr)
		if runOpts.Log.Runner.LogStderr {
			if runOpts.Log.Runner.LogStderrToStdout {
				cmd.Stderr = io.MultiWriter(sysout, stdoutWr)
			} else {
				cmd.Stderr = io.MultiWriter(syserr, stdoutWr)
			}
		} else {
			cmd.Stderr = stdoutWr
		}
	} else {
		if Server.Log.Runner.LogStderr {
			stderrOutput := syserr
			if runOpts.Log.Runner.LogStderrToStdout {
				stderrOutput = sysout
			}
			stderrWr = io.MultiWriter(stderrOutput, &stderr_buf)
		} else {
			stderrWr = &stderr_buf
		}
		cmd.Stderr = stderrWr
	}

	cmd.Stdout = stdoutWr

	err = cmd.Run()

	stderrStr := stderr_buf.String()
	stdoutStr := stdout_buf.String()
	
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

func parseAndRun(src string, registry map[string]Runner, req *http.Request, wr http.ResponseWriter) (string, error) {
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
			_, _, ret := errRun("create temporary file", err)
			return "", ret
		};defer os.RemoveAll(tmpDir)

		err = genLib(lang, req, tmpDir)
		if err != nil {
			_, _, ret := errRun("create temporary file", err)
			return "", ret
		}
		
		tmp := prepForLangsWithOddReqs(lang, tmpDir)

		code = formatCode(code, lang, tmp.Name(), tmpDir)

		api := BldApiStruct(req, wr)
		runOpts := RunOpts{
			Code: code,
			TmpDir: tmpDir,
			TmpFile: tmp,
			Req: req,
			Wr: wr,
			Lang: lang,
			Registry: registry,
			Api: api,
		}
		stdout, stderr, err := r.Run(runOpts)
		if err != nil {
			return "", fmt.Errorf("runner %s failed: %w; stderr=%s", lang, err, stderr)
		}

		stdout = formatSTD(lang, stdout)

		out.WriteString(stdout)
		i = end + 2
	}
	return out.String(), nil
}


func prepForLangsWithOddReqs(lang string, tmpDir string) *os.File {
	switch (lang) {
	case "go":
		modCont := []byte("module elh\n\ngo 1.20\n")
		err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), modCont, 0644)
		if err != nil {
			if Server.Log.Func != nil {
				errStr := fmt.Sprintf("err prepping for go: %v", err)
				Server.Log.Func(errStr)
			} else { fmt.Fprintf(syserr, "err prepping for go:  %v\n", err) }
			return nil
		}
		file, err := os.Create(filepath.Join(tmpDir, "main.go"))
		if err != nil {
			if Server.Log.Func != nil {
				errStr := fmt.Sprintf("%v", err)
				Server.Log.Func(errStr)
			} else { fmt.Fprintf(syserr, "%v\n", err) }
			return nil
		}
		return file
	default:
		fileName := fmt.Sprintf(filepath.Base(tmpDir))
		file, err := os.Create(filepath.Join(tmpDir, fileName))
		if err != nil {
			if Server.Log.Func != nil {
				errStr := fmt.Sprintf("%v", err)
				Server.Log.Func(errStr)
			} else { fmt.Fprintf(syserr, "%v\n", err) }
			return nil
		}
		return file
	}
}

func errRun(str string, err error) (string, string, error) {
	return "", "", fmt.Errorf("%s:  %w", str, err)
}

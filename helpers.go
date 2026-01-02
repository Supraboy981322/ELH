package ELH

import (
	"os"
	"fmt"
	"time"
	"errors"
	"net/http"
	"path/filepath"
)

// render src with specific registery
func RenderWithRegistry(src string, registry map[string]Runner, r *http.Request, wr http.ResponseWriter) (string, error) {
	return parseAndRun(src, registry, r, wr)
}

// wrapper that uses the DefaultRegistry.
func Render(src string, r *http.Request, wr http.ResponseWriter) (string, error) {
	return RenderWithRegistry(src, DefaultRegistry(), r, wr)
}

func MkReg(caller string, cmd string, args []string, timeout int, env []string) map[string]Runner {
	reg := map[string]Runner {
		caller: &ExternalRunner {
			CmdName: cmd,
			Args:    args,
			Timeout: time.Duration(timeout) * time.Second,
			Env:     env,
		},
	}
	return reg
}

func MkRegDefaults(cmd string, args []string) map[string]Runner {
	reg := MkReg(cmd, cmd, args, 5, os.Environ())
	return reg
}

func UseDefault(name string) *ExternalRunner {
	reg := DefaultRegistry()[name]
	if reg == nil {
		if Server.Log.Func != nil {
			Server.Log.Func("runner not in registry:  \n"+name)
			os.Exit(1)
		}
		fmt.Fprintf(syserr, "runner not in registry:  %s\n", name)
		os.Exit(1)
	}
	return reg.GetRunner()
}

func RenderFile(file string, r *http.Request, w http.ResponseWriter) ([]byte, error) {
	//get the extension of the requested file
	ext := filepath.Ext(file)	
	if ext == "" { //if there is no ext
		//check against list of ext which can
		//  have no ext in url
		for i := 0; i < len(suppNoExt); i++ {
			checkFile := fmt.Sprintf("%s%s", file, suppNoExt[i])
			_, err := os.Stat(checkFile)
			if err == nil { //if the file exists
				file = checkFile //assume it's the correct one
				ext = suppNoExt[i]
				break
			} else if !errors.Is(err, os.ErrNotExist) {
				return nil, errors.New("cannot check if file exists! Schrodinger's file"+err.Error())
			}
		}
	}

	fileByte, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.New("read file"+err.Error())
	}

	//if the file is elh, parse it
	if ext == ".elh" {
		fileStr := string(fileByte)
		result, err := Render(fileStr, r, w)
		if err != nil {
			return nil, errors.New("elh failed:  "+err.Error())
		}
		return []byte(result), nil
	} else {
		return fileByte, nil
	}

	return nil, errors.New("elh failed: uncaught err")
}

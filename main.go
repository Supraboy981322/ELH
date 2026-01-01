package ELH

import (
	"os"
	"fmt"
	"errors"
	"net/http"
	"path/filepath"
)

var (
	//for urls with no file format 
	suppNoExt = []string{
		".elh",
		".html",
	}
	Server ServerOpts
	sysout = os.Stdout
	syserr = os.Stderr
)

type (
	urlApi struct{
		Host string
		Path string
		Params map[string]string
		URI string
	}
	reqApi struct{
		Method string
		IsTLS bool
		URL urlApi
	}
	API struct {
		Req reqApi
	}
	RunOpts struct {
		Code string
		TmpDir string
		TmpFile *os.File
		Req *http.Request
		Wr http.ResponseWriter
		Lang string
		Registry map[string]Runner
		Api API
	}
	Runner interface {
		Run(runOpts RunOpts) (stdout string, stderr string, err error)
		GetRunner() (*ExternalRunner)
	}
	Logger struct {
		Func func(string)
		Runner struct {
			LogStderr bool
			LogStderrToStdout bool
		}
	}
	ServerOpts struct {
		Port int	
		Serve func() error
		WebDir string
		Log Logger
		ErrPage []byte
		RenderStderr bool
	}
)

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

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

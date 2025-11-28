package ELH

import (
	"os"
	"fmt"
	"errors"
	"net/http"
	"path/filepath"
)

const (
	port = "8080"	
)

var (
	//for urls with no file format 
	suppNoExt = []string{
		".elh",
		".html",
	}
)

type (
	Runner interface {
		Run(code string, tmp *os.File) (stdout string, stderr string, err error)
	}
)

func RenderFile(file string, r *http.Request) ([]byte, error) {
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
		result, err := Render(fileStr, r)
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

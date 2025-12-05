package ELH

import (
	"os"
	"fmt"
	"time"
	"bytes"
	"errors"
	"net/http"
	"path/filepath"
)

// render src with specific registery
func RenderWithRegistry(src string, registry map[string]Runner, r *http.Request) (string, error) {
	return parseAndRun(src, registry, r)
}

// wrapper that uses the DefaultRegistry.
func Render(src string, r *http.Request) (string, error) {
	return RenderWithRegistry(src, DefaultRegistry(), r)
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

func Serve(w http.ResponseWriter, r *http.Request) (string, error) {
	//get the requested file
	file := r.URL.Path
	if file == "/" {
		file = "index"
	} else if file[len(file)-1:] ==  "/" {
		file = fmt.Sprintf("%sindex", string(file[1:]))
		file, _ = checkIsDir(file)
	} else {
		file = file[1:]
		file, _ = checkIsDir(file)
	}

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
				return file, errors.New("cannot check if file exists! Schrodinger's file:  "+err.Error())
			}
		}
	}
	if fileExists(file) {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			return file, errors.New("read file:  "+err.Error())
		}
		fileStr := string(fileByte)
		var result string
		//if the file is elh, parse it
		if ext == ".elh" {
			result, err = Render(fileStr, r)
			if err != nil {
				return file, errors.New("elh failed; "+err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				fmt.Fprintf(w, "There appears to be an error in the `.elh` file %s", file)
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	} else {
		http.Error(w, "404 forbidden", http.StatusForbidden)
		file = "404 forbidden" 
	}
	return file, nil
}


//returns error so Go doesn't panic, but error
//  is ignored when call fn
//   (it's handled later, after fn call) 
func checkIsDir(file string) (string, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return file, err
	}
	if fileInfo.IsDir() {
		file = fmt.Sprintf("%s/index", file)
	}
	return file, nil
}

func ServeWithRegistry(w http.ResponseWriter, r *http.Request, registry map[string]Runner) (string, error) {
	//get the requested file
	file := r.URL.Path
	//set the system path of file
	if file == "/" {
		file = "index"
	} else if file[len(file)-1:] ==  "/" {
		file = fmt.Sprintf("%sindex", string(file[1:]))
		file = filepath.Join(WebDir, file)
		file, _ = checkIsDir(file)
	} else {
		file = file[1:]
		file = filepath.Join(WebDir, file)
		fmt.Printf("WebDir=%s ; file=%s\n", WebDir, file)
		file, _ = checkIsDir(file)
	}

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
				return file, errors.New("cannot check if file exists! Schrodinger's file:  "+err.Error())
			}
		}
	}
	if fileExists(file) {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			return file, errors.New("read file:  "+err.Error())
		}
		fileStr := string(fileByte)
		var result string
		//if the file is elh, parse it
		if ext == ".elh" {
			result, err = RenderWithRegistry(fileStr, registry, r)
			if err != nil {
				return file, errors.New("elh failed; "+err.Error())
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				fmt.Fprintf(w, "There appears to be an error in the `.elh` file %s", file)
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	} else {
		http.Error(w, "404 forbidden", http.StatusForbidden)
		file = "404 forbidden" 
	}
	return file, nil
}

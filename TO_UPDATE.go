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
			result, err = Render(fileStr, r, w)
			if err != nil {
				var fileStr string
				if Server.ErrPage != nil {
					fileStr, err = Render(string(Server.ErrPage), r, w) 
					if err != nil {
						http.Error(w, "500 server err", 500)
						return "500", err
					}
				} else { fileStr = "500 server err" }
				http.Error(w, fileStr, 500)
				file = "500 server err" 
				return file, errors.New("elh failed; "+err.Error())
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	} else {
		var err error
		var fileStr string
		if Server.ErrPage != nil {
			fileStr, err = Render(string(Server.ErrPage), r, w) 
			if err != nil {
				http.Error(w, "500 server err", 500)
				return "500", err
			}
		} else { fileStr = "404 forbidden" }
		http.Error(w, fileStr, 404)
		file = "404 forbidden" 
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
		file = filepath.Join(Server.WebDir, file)
		file, _ = checkIsDir(file)
	} else {
		file = file[1:]
		file = filepath.Join(Server.WebDir, file)
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
			result, err = RenderWithRegistry(fileStr, registry, r, w)
			if err != nil {
				var fileStr string
				if Server.ErrPage != nil {
					fileStr, err = Render(string(Server.ErrPage), r, w) 	
					if err != nil {
						http.Error(w, "500 server err", 500)
						return "500", err
					}
				} else { fileStr = "500 server err" }
				http.Error(w, fileStr, 500)
				file = "500 server err" 
				return file, errors.New("elh failed; "+err.Error())
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	} else {
		var err error
		var fileStr string
		if Server.ErrPage != nil {
			fileStr, err = Render(string(Server.ErrPage), r, w) 
			if err != nil {
				http.Error(w, "500 server err", 500)
				return "500", err
			}
		} else { fileStr = "404 forbidden" }
		http.Error(w, fileStr, 404)
		file = "404 forbidden" 
	}
	return file, nil
}

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

func (server *ServerOpts) HttpHandler(w http.ResponseWriter, r *http.Request) {
	if server.Registry == nil {
		server.Registry = DefaultRegistry()
	}
	var resp string
//	if server.WebDir == "" { server.WebDir, _ = os.Getwd() }
	if server.WebDir == "" { server.WebDir = "." }
	//get the requested file
	file := filepath.Join(server.WebDir, r.URL.Path)
	if file == server.WebDir {
		file = filepath.Join(server.WebDir, "index")
	} else {
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
				http.Error(w, "cannot check if file exists! Schrodinger's file:  "+err.Error(), http.StatusInternalServerError)
				resp = "500; Schrodinger's file"
			}
		}
	}
	if fileExists(file) {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			http.Error(w, "read file:  "+err.Error(), http.StatusInternalServerError)
			resp = "500; err reading file"
		}
		fileStr := string(fileByte)
		var result string
		//if the file is elh, parse it
		if ext == ".elh" {
			result, err = RenderWithRegistry(fileStr, server.Registry, r, w)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				fmt.Fprintf(w, "There appears to be an error in the `.elh` file %s", file)
				if server.Log.Func != nil {
					str := fmt.Sprintf("[err]: RenderWithRegistry() %v", err)
					server.Log.Func(str)
				}
				resp = "500; problem with elh file"
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
		if server.ErrPage != nil {
			fileStr, err = Render(string(server.ErrPage), r, w)
			if err != nil {
				http.Error(w, "500 server err", 500)
				return
			}
		} else { fileStr = "404 forbidden" }
		http.Error(w, fileStr, 404)
		file = "404 forbidden" 
	}

	//colorize response log string
	if resp == "" { resp = "\033[32m"+file+"\033[0m"
	} else { resp = "\033[31m"+resp+"\033[0m" }

	//colorize file log string
	file = "\033[35m"+file+"\033[0m"

	//check if logger is set
	if server.Log.Func != nil {
		//build string
		logStr := "\033[1m[req]:\033[0m "
		logStr += file+" | "
		logStr += "\033[1m[resp]:\033[0m "+resp
		//log it
		server.Log.Func(logStr)
	}
}

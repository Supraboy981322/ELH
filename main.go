package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"os"
	"errors"
	"bytes"
	"time"
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

type Runner interface {
	Run(code string, tmp *os.File) (stdout string, stderr string, err error)
}


func main() {
	http.HandleFunc("/", httpHandler)
	fmt.Printf("listening on port:  %s\n", port)
	errOut("http", http.ListenAndServe(":"+port, nil))
}


func httpHandler(w http.ResponseWriter, r *http.Request) {
	//check method type
	switch r.Method {
	case http.MethodGet: //run fn for GET
		fmt.Printf("GET:  %s\n", r.URL.Path)
		getHandler(w, r)
	default: //only GET is supported at the moment
		fmt.Printf("attempted bad method:  %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	//get the requested file
	file := r.URL.Path
	if file == "/" {
		file = "index"
	} else if file[len(file)-1:] ==  "/" {
		file = fmt.Sprintf("%sindex", string(file[1:]))
	} else {
		file = file[1:]
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
				errOut("cannot check if file exists! Schrodinger's file", err)
			}
		}
	}
	if fileExists(file) {
		fileByte, err := os.ReadFile(file)
		if err != nil {
			errOut("read file", err)
		}
		fileStr := string(fileByte)
		var result string
		//if the file is elh, parse it
		if ext == ".elh" {
			result, err = Render(fileStr)
			if err != nil {
				errOut("elh failed;", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				fmt.Fprintf(w, "There appears to be an error in the `.elh` file %s", file)
			}
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprintln(w, result)
		} else {
			fileReader := bytes.NewReader(fileByte)
			http.ServeContent(w, r, file, time.Now(), fileReader)
		}
	
		//log request
		fmt.Printf("\nreq:  %s\n", file)	
	} else {
		http.Error(w, "404 forbidden", http.StatusForbidden)
	}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}


//for errors
func errOut(str string, err error) {
	errStr := fmt.Sprintf("\n\n%s:\n%v\n", str, err)
	errVal := errors.New(errStr)
	log.Println(errVal)
}

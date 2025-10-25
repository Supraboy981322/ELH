package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"os"
	"errors"
)

const (
	port = "8080"
)

var (
	suppNoExt = []string{
		".html",
		".egp",
	}
)

func main() {
	http.HandleFunc("/", httpHandler)

	panic(http.ListenAndServe(":"+port, nil))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		fmt.Printf("GET:  %s\n", r.URL.Path)
		getHandler(w, r)
	default:
		fmt.Printf("attempted bad method:  %s\n", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Path
	if file ==  "/" {
		file = "index"
	}
	ext := filepath.Ext(file)
	fmt.Println(ext)
	if ext == "" {
		for i := 0; i < len(suppNoExt); i++ {
			checkFile := fmt.Sprintf("%s%s", file, suppNoExt[i])
			_, err := os.Stat(checkFile)
			if err == nil {
				file = checkFile
				ext = suppNoExt[i]
				break
			} else if !errors.Is(err, os.ErrNotExist) {
				errOut("cannot check if file exists! Schrodinger's file", err)
			}
		}
	}

	if ext == ".egp" {
		
	}

	fmt.Printf("\nreq:  %s\n", file)

	fmt.Fprintf(w, "foo")
}

func errOut(str string, err error) {
	errStr := fmt.Sprintf("\n\n%s:\n%v\n", str, err)
	errVal := errors.New(errStr)
	panic(errVal)
}

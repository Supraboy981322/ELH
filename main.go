package main

import (
	"fmt"
	"net/http"
//	"os"
)

const (
	port = "8080"
)

func main() {
	http.HandleFunc("/", httpHandler)

	panic(http.ListenAndServe(":"+port, nil))
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\nreq:  %s\n", r.URL.Path)
	
	fmt.Fprintf(w, "foo")
}

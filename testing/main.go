package main

import (
	"os"
	"time"
	"net/http"
	"github.com/charmbracelet/log"
	elh "github.com/Supraboy981322/ELH"
	"github.com/Supraboy981322/gomn"
)
var gomnMap gomn.Map

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		renderFromFileName("test.elh", w, r)
		renderFromRegistry("test.elh", w, r)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}



/* * * * * * * * * * * * *
 * render from file name *
 * * * * * * * * * * * * */
func renderFromFileName(file string, w http.ResponseWriter, r *http.Request) {
	res, err := elh.RenderFile("test.elh", r)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(res)
}



/* * * * * * * * * * * * * * * *
 * render with custom registry *
 * * * * * * * * * * * * * * * */
func renderFromRegistry(file string, w http.ResponseWriter, r *http.Request) {
	//read file 
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	//make custom map
  //	registry := elh.MkRegDefaults("bash", []string{})
	//  registry := elh.MkReg("bash", "bash', []string{}, 5, os.Environ())
	registry := map[string]elh.Runner{
		"bash": &elh.ExternalRunner{
			CmdName: "bash",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
	}

	res, err := elh.RenderWithRegistry(string(fileBytes), registry, r)
	if err != nil {
		log.Fatal(err)
	}
	w.Write([]byte(res))
}

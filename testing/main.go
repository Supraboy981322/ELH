package main

import (
	"os"
	"time"
	"net/http"
	"github.com/charmbracelet/log"
	elh "github.com/Supraboy981322/ELH"
	"github.com/Supraboy981322/gomn"
)
var (
	gomnMap gomn.Map
	port = "8080"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Infof("[req]:  %s", r.RemoteAddr)
		renderFromFileName("index.elh", w, r)
//		renderFromRegistry("test.elh", w, r)
	})

	log.Infof("listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}



/* * * * * * * * * * * * *
 * render from file name *
 * * * * * * * * * * * * */
func renderFromFileName(file string, w http.ResponseWriter, r *http.Request) {
	res, err := elh.RenderFile(file, r)
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
	//	registry := elh.MkReg("bash", "bash", []string{}, 5, os.Environ())
  //	registry := elh.MkRegDefaults("bash", []string{})
	//  registry := elh.MkReg("bash", "bash', []string{}, 5, os.Environ())
	registry := map[string]elh.Runner{
		"bash": &elh.ExternalRunner{
			CmdName: "bash",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"md": &elh.ExternalRunner{
			CmdName: "marked",
			Args:    []string{"-i"},
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

package ELH

import (
	"os"
//	"strconv"
	"strings"
	"net/http"
//	"encoding/json"
	"path/filepath"
)


func formatCode(code string, lang string, tmpName string, tmpDir string) string { 
	impArray := getImpsBetween(code, "<??imps", "??>")
	code = stripImps(code)
	var imps string
	switch (lang) {
	case "java":
		fileName := strings.ReplaceAll(tmpName, tmpDir + "/", "")
		class := fileName
		code = "public class "+ class + " {\n" + code + "\n}\n"
	case "go":
		var head string
		if impArray[0] != ""  {
			for i := 0; i < len(impArray); i++ {
				imps = strings.TrimSpace(impArray[i]) + "\n"
			}
			head = "package main\nimport (\n" + imps + "\n)\n"
		} else {
			head = "package main\n"
		}
		code = head + code
		sysout.WriteString(code)
	case "php":
		code = "<?php\n" + code + "\n?>"
	case "py":
		if impArray[0] != "" {
			var head string
			for i := 0; i < len(impArray); i++ {
				head += "import "
				head += strings.TrimSpace(impArray[i])
				head += "\n"
			}
			code = head + "\n" + code
		}
	case "basic":
		code = code + "\nQUIT"
	case "vim":
		if impArray[0] != "" {
			var head string
			for i := 0; i < len(impArray); i++ {
				imp := strings.TrimSpace(impArray[i])
				head += "source "
				if imp != "elh" {
					head += imp
				} else {
					head += filepath.Join(tmpDir, "elhLib.vim")
				}
				head += "\n"
			}
			code = head + "\n" + code
		}
		code = code + "\nqall!"
	case "bash":
		head := `
source `+filepath.Join(tmpDir, "elhLib.bash")+"\n"
		code = head + code
	default:
	}
	return code
}


func formatSTD(lang string, stdout string) string {
	res := stdout
	switch lang {
	case "basic":
		stdLi := strings.Split(stdout, "\n")
		stdLi = stdLi[3:]
		res = strings.Join(stdLi, "\n") 
	case "vim":
	default:
	}
	if len(res) > 1 { 
		if res[len(res)-1] == '\n' { res = res[:len(res)-2] }
	}
	return res
}

func stripImps(code string) string {
	start := strings.Index(code, "??>")
	if start == -1 {
		return code
	}
	start += 3
	return code[start:]
}

func formatHeaders(r *http.Request, lang string) string {
	var headArr string
	switch (lang) {
	case "vim":
		headArr = "{"
		for name, values := range r.Header {
			for _, value := range values {
				headArr += "'" + name + "':'" + value + "',"
			}
		}
		headArr += "}"
	case "bash":
/*		jsonHeaders, err := json.Marshal(r.Header)
		if err != nil {
			return err.Error()
		}
		headArr = string(jsonHeaders)*/
	default:
	}
	return headArr
}

func fmtdHead(r *http.Request, lang string) string {
	var headArr string
	switch (lang){
	case "vim":
		headArr = "{"
		for name, values := range r.Header {
			for _, value := range values {
				headArr += " '" + name + "': '" + value + "',"
			}
		}
		headArr += " }"
	default:
	}
	return headArr
}

func genLib(lang string, r *http.Request, tmpDir string) error {
//	headers := fmtdHead(r, lang)
	switch (lang) {
	case "vim":
		libCont := []byte(`
let s:Headers = " + headers + "
let s:Params = { 'TODO': 'TODO' }
let elh = { 'Headers': s:Headers, 'Params': s:Params, }
`)
		libName := filepath.Join(tmpDir, "elhLib.vim")
		err := os.WriteFile(libName, libCont, 0644)
		if err != nil {
			return err
		}
	 case "bash":
		libCont := []byte(`
elhlog() (
	set -euo pipefail
	printf "${@}" 1>&2
	printf "\n" 1>&2
)
`)
		libName := filepath.Join(tmpDir, "elhLib.bash")
		err := os.WriteFile(libName, libCont, 0644)
		if err != nil {
			return err
		}
  default:
	}
	return nil
}

func getImpsBetween(code string, start string, end string) []string {
	res := []string{""}
	
	star := strings.Index(code, start)
	if star == -1 {
//		res[0] = code
		return res
	}
	star += len (start)

	en := strings.Index(code[star:], end)
	if en == -1 {
		res[0] = code
		return res
	}

	res[0] = code[star : star+en]
	res = strings.Split(res[0], " ; ")
	return res
}

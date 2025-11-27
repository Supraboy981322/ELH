package main

import (
	"os"
	"strings"
	"net/http"
	"path/filepath"
	"strconv"
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
		os.Stdout.WriteString(code)
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
	headers := fmtdHead(r, lang)
	switch (lang) {
	case "vim":
		libArr := []string{
			"let s:Headers = " + headers + "\n",
			"let s:Params = { 'TODO': 'TODO' }\n",
			"let elh = {",
			" 'Headers': s:Headers,",
			" 'Params': s:Params,",
			" }\n",
		}
		var libCont string
		for i := 0; i < len(libArr); i++ {
			libCont += libArr[i]
		}
		libName := filepath.Join(tmpDir, "elhLib.vim")
		err := os.WriteFile(libName, []byte(libCont), 0644)
		if err != nil {
			return err
		}
	default:
	}
	return nil
}

func getImpsBetween(code string, start string, end string) []string {
	res := []string{""}
	
	os.Stdout.WriteString(strconv.Itoa(1))
	star := strings.Index(code, start)
	if star == -1 {
//		res[0] = code
		return res
	}
	star += len (start)

	os.Stdout.WriteString(strconv.Itoa(2))
	en := strings.Index(code[star:], end)
	if en == -1 {
		res[0] = code
		return res
	}

	os.Stdout.WriteString(strconv.Itoa(3))
	res[0] = code[star : star+en]
	res = strings.Split(res[0], " ; ")
	return res
}

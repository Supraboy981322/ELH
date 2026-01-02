package ELH

import (
	"os"
	"strings"
	"net/http"
)

func BldApiStruct(req *http.Request, wr http.ResponseWriter) API {
	params := map[string]string{}
	for k, _ := range req.URL.Query() {
		params[k] = req.URL.Query().Get(k)
	}
	api := API{
		Req: reqApi{
			Method: req.Method,
			IsTLS: false,
			URL: urlApi{
				Host: req.Host,
				URI: req.RequestURI,
				Path: req.URL.Path,
				Params: params,
			},
		},
	}
	
	if req.TLS != nil { api.Req.IsTLS = true }

	return api
}

func ToArgs(src string, by string) []string {
	if len(src) <= 1 { return []string{} }
	if src[len(src)-1] == ' ' { src = src[:len(src)] }
	if len(src) <= 1 { return []string{} }
	type argParser struct{
		quot bool
		esc bool
		pos int
		mem string
		out []string
		in []string
	};p := argParser{
		quot: false,
		esc: false,
		pos: 0,
		mem: "",
		out: []string{},
		in: strings.Split(src, ""),
	}

	var foo func(p argParser) []string
	foo = func(p argParser) []string {
		if p.pos >= len(p.in) {
			p.out = append(p.out, p.mem)
			return p.out
		}
		cur := p.in[p.pos]
		if p.esc || p.quot {
			if p.quot && (cur == `"` || cur == `'`) {
				p.pos++ ; p.quot = false ; return foo(p)
			}
			p.mem += cur
			p.pos++
			if p.esc { p.esc = false }
			return foo(p)
		}
		switch cur {
     case `\`: p.esc = true
	   case `"`, `'`: p.quot = true
		 case by:
			p.out = append(p.out, p.mem)
			p.mem = ""
		 default: p.mem += cur
		}
		p.pos++
		return foo(p)
	}
	return foo(p)
}

func (opts *RunOpts) Unwrap() (string, string, *os.File, *http.Request, http.ResponseWriter) {
	code := opts.Code
	tmpDir := opts.TmpDir
	tmpFile := opts.TmpFile
	req := opts.Req
	wr := opts.Wr
	return code, tmpDir, tmpFile, req, wr
}


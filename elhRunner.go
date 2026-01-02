package ELH

import (
	"os"
	"fmt"
	"strconv"
	"unicode"
	"strings"
)

var ErrInvalidPath = fmt.Errorf("invalid path")

func elhRunner(opts RunOpts) (string, string, error) {
	var res any
	if args := splitLines(opts.Code, []string{";"}); len(args) > 1 {
		var stdout, stderr string
		for _, l := range args {
			o := opts ; o.Code = l
			stdoutT, stderrT, err := elhRunner(o)
			stdout += stdoutT ; stderr += stderrT
			if err != nil { return stdout, stderr, err }
		}
		return stdout, stderr, nil
	}
	args := ToArgs(opts.Code, ".")
	if len(args) < 1 {
		err := fmt.Errorf("no args provided")
		return "", "no arg", err
	}
	if len(opts.Code) >= 2 {
		if strings.HasPrefix(opts.Code, "//") {
			return "", "", nil
		}
	}
	type pa struct{
		esc bool
		pos int
		memInt int
		mem string
		mem2 []string
		in []string
		out string
	}
	invArg := func(a string, r string) (string, string, error) {
		err := fmt.Errorf("called '%s', but %s", a, r)
		return "", r, err
	}
	invPath := func(p string) (string, string, error){
		err := fmt.Errorf("invalid path: '%s'", p)
		return "500 server err", "invalid path", err
	}
	var getString func(p pa) (string, []string)
	getString = func(p pa) (string, []string) {
		if p.pos >= len(p.in) { return p.out, p.mem2 }
		switch p.in[p.pos] {
		 case "(":
			if !p.esc {
				p.out = p.mem
				p.mem = ""
				p.memInt = 0
			}
		 case ",":
			if !p.esc { p.memInt++ }
		 case " ":
			if !p.esc {
				p.pos++ ; return getString(p)
			} else { p.mem += p.in[p.pos] }
		 case ")":
			if !p.esc {
				p.mem2 = append(p.mem2, p.mem)
			}
		 case "'": p.esc = !p.esc
		 default:
			p.mem += p.in[p.pos]
		}
		p.pos++
		return getString(p)
	}
	var getIndex func(p pa) (string, int)
	getIndex = func(p pa) (string, int) {
		if p.pos >= len(p.in) { return p.out, p.memInt }
		switch p.in[p.pos] {
		 case "[": p.out += p.mem;p.mem = ""
		 case "]":
			var err error
			p.memInt, err = strconv.Atoi(p.mem)
			if err != nil { return "", 0 }
			p.pos++
		 default: p.mem += p.in[p.pos]
		};p.pos++
		return getIndex(p)
	}
	for i, a := range args { args[i] = strings.TrimSpace(a) }
	if strings.Contains(args[0], "(") {
		cmd, fi := getString(pa{ in: strings.Split(args[0], "") })
		args = append(append([]string{cmd}, fi...), args[1:]...)
	}
	if unicode.IsUpper(rune(args[0][0])) {
		sT, err := readStruct(opts, strings.Join(args, "."))
		if err != nil { return "", err.Error(), err } 
		res = fmt.Sprintf("%v", sT.String())
	} else {
		switch args[0] {
		 case "req":
			if len(args) < 2 { return invArg(args[0], "path ended too soon") }
			switch args[1] {
			 case "header":
				if len(args) < 3 { return invPath(strings.Join(args, ".")) }
				header, i := getIndex(pa{ in: strings.Split(args[2], "") })
				if header == "" { return invPath(args[2]) }
				res = opts.Req.Header[header][i]
			}
		 case "dump_file":
			if len(args) < 2 { return invArg(args[0], "file not provided") }
			resT, err := os.ReadFile(args[1])
			res = string(resT)
			if err != nil { return "500", "failed to read file", err }
		 case "render_file":
			var useReg bool
			var file string
			if len(args) < 2 { return invArg(args[0], "file not provided") }
			file = args[1]
			fiB, err := os.ReadFile(args[1])
			fi := string(fiB)
			if err != nil { return "500", "failed to read file", err }
			if len(args) >= 3 { useReg = strings.ToLower(args[1]) == "true" }
			var resT string 
			if useReg {
				resT, err = RenderWithRegistry(fi, opts.Registry, opts.Req, opts.Wr)
				if err != nil { return "500", "failed to read file", err }
			} else {
				var resB []byte
				resB, err = RenderFile(file, opts.Req, opts.Wr)
				resT = string(resB)
			}
			res = string(resT)
			if err != nil { return "500", "failed to read file", err }
		}
	}
	if res == nil { return invPath(strings.Join(args, ".")) }
	resF := fmt.Sprintf("%v", res)
	return resF, "", nil
}

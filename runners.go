package ELH

import (
	"io"
	"os"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"github.com/Shopify/go-lua"
	"github.com/gomarkdown/markdown"
	"github.com/Supraboy981322/gomn"
)

func mdRunner(opts RunOpts) (string, string, error) {
	if opts.Code == "" {
		err := fmt.Errorf("markdown is empty")
		return "err, see server logs", "no input", err
	}

	//render HTML as md 
	res := markdown.ToHTML([]byte(opts.Code), nil, nil)

	return string(res), "", nil
}

func luaRunner(opts RunOpts) (string, string, error) {
	outR, outW, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(syserr, "lua failed:  %v", err)
		return "", "", err
	}

	os.Stdout, os.Stderr = outW, errW

	l := lua.NewState()
	lua.OpenLibraries(l)
	
	if err := lua.DoString(l, opts.Code); err != nil {
		os.Stdout, os.Stderr = sysout, syserr
		fmt.Fprintf(syserr, "lua failed:  %v", err)
		return "", "", err
	}
	os.Stdout, os.Stderr = sysout, syserr
	errW.Close()
	outW.Close()

	outF, err := io.ReadAll(outR)
	if err != nil {
		outR.Close()
		fmt.Fprintf(syserr, "lua failed:  %v", err)
		return "", "", err
	};outR.Close()
	errF, err := io.ReadAll(errR)
	if err != nil {
		errR.Close()
		fmt.Fprintf(syserr, "lua failed:  %v", err)
		return string(outF), "", err
	};errR.Close()

	return string(outF), string(errF), err
}

func gomnParser(opts RunOpts) (string, string, error) {
	var file string
	var path []string
	args := ToArgs(opts.Code, " ")

	var taken []int
	for i, arg := range args {
		if !slices.Contains(taken, i) {
			switch arg {
			 case "-f":
				if len(args) <= i+1 {
					err := fmt.Errorf("called %s arg, but no value provided", arg)
					return "", "", err
				}
				file = args[i+1]
				taken = append(taken, i+1)
			 default:
				path = append(path, arg)
			}
		}
	}

	if file == "-" { file = "" }
	inB, err := os.ReadFile(file)
	if err != nil { return "", "failed to read file", err }
	in := string(inB)

	GOMN, err := gomn.Parse(in)
	if err != nil {
		return "", "failed to parse gomn", fmt.Errorf("\r%v", err)
	}

	var resA any
	cur := GOMN
	for _, p := range path {
		if n, ok := cur[p].(gomn.Map); ok {
				cur = n
		} else { resA = cur[p] }
	}

	res := fmt.Sprint(resA)

	return res, "", nil
}

func elhRunner(opts RunOpts) (string, string, error) {
	var res any
	args := ToArgs(opts.Code, ".")
	if len(args) < 1 {
		err := fmt.Errorf("no args provided")
		return "", "no arg", err
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
	if strings.Contains(args[0], "(") {
		cmd, fi := getString(pa{ in: strings.Split(args[0], "") })
		args = append(append([]string{cmd}, fi...), args[1:]...)
	}
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
		} else {
			var resB []byte
			resB, err = RenderFile(file, opts.Req, opts.Wr)
			resT = string(resB)
		}
		res = string(resT)
		if err != nil { return "500", "failed to read file", err }
	}
	if res == nil { return invPath(strings.Join(args, ".")) }
	resF := fmt.Sprintf("%v", res)
	return resF, "", nil
}

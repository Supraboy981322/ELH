package ELH

import (
	"io"
	"os"
	"fmt"
	"slices"
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


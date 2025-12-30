package ELH

import (
	"io"
	"os"
	"fmt"
	"github.com/Shopify/go-lua"
	"github.com/gomarkdown/markdown"
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
	sysout, syserr := os.Stdout, os.Stderr

	outR, outW, err := os.Pipe()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "lua failed:  %v", err)
		return "", "", err
	}

	os.Stdout, os.Stderr = outW, errW

	l := lua.NewState()
	lua.OpenLibraries(l)
	
	if err := lua.DoString(l, opts.Code); err != nil {
		os.Stdout, os.Stderr= sysout, syserr
		fmt.Fprintf(os.Stderr, "lua failed:  %v", err)
		return "", "", err
	}
	os.Stdout, os.Stderr = sysout, syserr
	errW.Close()
	outW.Close()

	outF, err := io.ReadAll(outR)
	if err != nil {
		outR.Close()
		fmt.Fprintf(os.Stderr, "lua failed:  %v", err)
		return "", "", err
	};outR.Close()
	errF, err := io.ReadAll(errR)
	if err != nil {
		errR.Close()
		fmt.Fprintf(os.Stderr, "lua failed:  %v", err)
		return string(outF), "", err
	};errR.Close()

	return string(outF), string(errF), err
}

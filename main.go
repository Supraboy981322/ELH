package ELH

import (
	"os"
	"net/http"
)

var (
	//for urls with no file format 
	suppNoExt = []string{
		".elh",
		".html",
	}
	Server ServerOpts
	sysout = os.Stdout
	syserr = os.Stderr
)

type (
	urlApi struct{
		Host string
		Path string
		Params map[string]string
		URI string
	}
	reqApi struct{
		Method string
		IsTLS bool
		URL urlApi
	}
	API struct {
		Req reqApi
	}
	RunOpts struct {
		Code string
		TmpDir string
		TmpFile *os.File
		Req *http.Request
		Wr http.ResponseWriter
		Lang string
		Registry map[string]Runner
		Api API
		Log Logger
		Server ServerOpts
	}
	Runner interface {
		Run(runOpts RunOpts) (stdout string, stderr string, err error)
		GetRunner() (*ExternalRunner)
	}
	Logger struct {
		Func func(string)
		Runner struct {
			LogStderr bool
			LogStderrToStdout bool
		}
	}
	ServerOpts struct {
		Port int	
		Serve func() error
		WebDir string
		Log Logger
		ErrPage []byte
		RenderStderr bool
		Registry map[string]Runner
	}
)

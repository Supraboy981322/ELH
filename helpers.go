package ELH

import (
	"os"
	"fmt"
	"time"
	"net/http"
)

// render src with specific registery
func RenderWithRegistry(src string, registry map[string]Runner, r *http.Request, wr http.ResponseWriter) (string, error) {
	return parseAndRun(src, registry, r, wr)
}

// wrapper that uses the DefaultRegistry.
func Render(src string, r *http.Request, wr http.ResponseWriter) (string, error) {
	return RenderWithRegistry(src, DefaultRegistry(), r, wr)
}

func MkReg(caller string, cmd string, args []string, timeout int, env []string) map[string]Runner {
	reg := map[string]Runner {
		caller: &ExternalRunner {
			CmdName: cmd,
			Args:    args,
			Timeout: time.Duration(timeout) * time.Second,
			Env:     env,
		},
	}
	return reg
}

func MkRegDefaults(cmd string, args []string) map[string]Runner {
	reg := MkReg(cmd, cmd, args, 5, os.Environ())
	return reg
}

func UseDefault(name string) *ExternalRunner {
	reg := DefaultRegistry()[name]
	if reg == nil {
		if Server.Log.Func != nil {
			Server.Log.Func("runner not in registry:  \n"+name)
			os.Exit(1)
		}
		fmt.Fprintf(syserr, "runner not in registry:  %s\n", name)
		os.Exit(1)
	}
	return reg.GetRunner()
}

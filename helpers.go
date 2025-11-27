package elh

import (
	"os"
	"time"
	"net/http"
)

// render src with specific registery
func RenderWithRegistry(src string, registry map[string]Runner, r *http.Request) (string, error) {
	return parseAndRun(src, registry, r)
}

// wrapper that uses the DefaultRegistry.
func Render(src string, r *http.Request) (string, error) {
	return RenderWithRegistry(src, DefaultRegistry(), r)
}

func MkReg(caller string, cmd string, args []string, timeout int, env []string) map[string]Runner {
	reg := map[string]Runner {
		caller: &ExternalRunner{
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

package main

import (
	"os"
	"time"
)

func DefaultRegistry() map[string]Runner {
	return map[string]Runner{
		"py": &ExternalRunner{
			CmdName: "python3",
			Args:    []string{"-u"}, // unbuffered stdout
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"js": &ExternalRunner{
			CmdName: "bun",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"bash": &ExternalRunner{
			CmdName: "bash",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"go": &ExternalRunner{
			CmdName: "go",
			Args: []string{}
			Timeout: 5 *time.Second,
			Env: os.Environ(),
		}
	}
}

// render src with specific registery
func RenderWithRegistry(src string, registry map[string]Runner) (string, error) {
	return parseAndRun(src, registry)
}

// wrapper that uses the DefaultRegistry.
func Render(src string) (string, error) {
	return RenderWithRegistry(src, DefaultRegistry())
}


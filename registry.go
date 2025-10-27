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
		"java": &ExternalRunner{
			CmdName: "java",
			Args:    []string{"--source", "23"},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"lua": &ExternalRunner{
			CmdName: "lua",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"go": &ExternalRunner{
			CmdName: "go",
			Args:    []string{"run"},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"bf": &ExternalRunner{
			CmdName: "beef",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"r": &ExternalRunner{
			CmdName: "Rscript",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"ruby": &ExternalRunner{
			CmdName: "ruby",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"php": &ExternalRunner{
			CmdName: "php",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"perl": &ExternalRunner{
			CmdName: "perl",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"md": &ExternalRunner{
			CmdName: "marked",
			Args:    []string{"-i"},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
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


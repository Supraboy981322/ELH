package ELH

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
			Func: jsRunner,
			Timeout: 5 * time.Second,
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
			Timeout: 5 * time.Second,
			Func: luaRunner,
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
		"basic": &ExternalRunner{
			CmdName: "bwbasic",
			Args:    []string{},
			Timeout: 5 * time.Second,
			Env:     os.Environ(),
		},
		"vim": &ExternalRunner{
			Timeout: 5 * time.Second,
			Func: vimRunner,
		},
		"md": &ExternalRunner{
			Timeout: 5 * time.Second,
			Func: mdRunner,
		},
		"gp": &ExternalRunner{
			Timeout: 5 * time.Second,
			Func: gomnParser,
		},
		"elh": &ExternalRunner{
			Timeout: 5 * time.Second,
			Func: elhRunner,
		},
	}
}


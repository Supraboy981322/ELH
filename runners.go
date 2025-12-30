package ELH

import (
	"os"
	"fmt"
	"github.com/gomarkdown/markdown"
)

func mdRunner(code string, tmp *os.File) (string, string, error) {
	if code == "" {
		err := fmt.Errorf("markdown is empty")
		return "err, see server logs", "no input", err
	}

	//render HTML as md 
	res := markdown.ToHTML([]byte(code), nil, nil)

	return string(res), "", nil
}

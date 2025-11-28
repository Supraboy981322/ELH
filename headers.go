package ELH

var (
	libs = map[string]string{
		//python libs
		"py":    ``,

		//js libs
		"js":    ``,

		//bash libs
		"bash":  `#!/usr/bin/env bash

headers='{"foo":"bar"}'

headers() {
  case "$1" in
    "get")
      if [[ "$2" == "" ]]; then
        printf "no header provided"
      fi
      ;;
    *)
      printf "$headers" | jq
      ;;
  esac
}

linesToHTML() {
  local stdin=$(< /dev/stdin) 
  printf "${stdin}\n" \
    | sed "s|^|<${1:-p}>|g" \
    | sed "s|$|</${1:-p}>|g"
}

linkFiles() {
  for file in $(ls); do
    printf "<a href=\"${file}\">${file}</a>\n"
  done
}`,

		//java libs
		"java":  ``,
		
		//lua libs
		"lua":   ``,

		//go libs
		"go":    ``,

		//brainfuck libs
		"bf":    ``,

		//r libs
		"r":     ``,

		//ruby libs
		"ruby":  ``,

		//php libs
		"php":   ``,

		//perl libs
		"perl":  ``,

		//basic libs
		"basic": ``,

		//vim libs
		"vim":   ``,
	}
)

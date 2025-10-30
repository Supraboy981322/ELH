package main

const (
	vimHead = "com! -nargs=* Import call Import(<f-args>)\nfunction! Import(lib)\n    exec \"source \" . a:lib . \".vim\"\nendfunction\n"
)

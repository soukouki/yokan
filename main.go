package main

import (
	"os"
	"yokan/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}

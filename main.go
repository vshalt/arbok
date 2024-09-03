package main

import (
	"os"

	"github.com/vshalt/arbok/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}

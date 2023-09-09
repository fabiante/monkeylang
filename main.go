package main

import (
	"github.com/fabiante/monkeylang/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}

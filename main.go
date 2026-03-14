package main

import (
	"fmt"
	"os"
)

var (
	name    = "ctagls"
	version = "да"
)

func main() {
	lspServer := NewLspServer(os.Stdin, os.Stdout)
	if err := lspServer.Handle(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

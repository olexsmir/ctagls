package main

import (
	"fmt"
	"log/slog"
	"os"
)

var (
	name    = "ctagls"
	version = "да"
)

func main() {
	logFile, err := os.OpenFile("ctagls.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	slog.SetDefault(slog.New(slog.NewTextHandler(
		logFile,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)))

	lspServer := NewLspServer(os.Stdin, os.Stdout)
	if err := lspServer.Handle(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

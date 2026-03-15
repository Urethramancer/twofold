package main

import (
	"os"

	"golang.org/x/term"
)

// isTerminal returns true when stdout is a terminal
func isTerminal() bool {
	fd := int(os.Stdout.Fd())
	return term.IsTerminal(fd)
}

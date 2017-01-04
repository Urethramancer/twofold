package main

const (
	program = "twofold"
)

// Version is filled in from Git tags.
var Version = "undefined"

func printVersion() {
	pr("%s %s", program, Version)
}

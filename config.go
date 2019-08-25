package main

const (
	program = "twofold"
)

// version is filled in from Git tags.
var version = "undefined"

func printVersion() {
	pr("%s %s", program, version)
}

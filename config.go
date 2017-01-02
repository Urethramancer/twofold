package main

const (
	program = "twofold"
)

var Version = "undefined"

func printVersion() {
	pr("%s %s", program, Version)
}

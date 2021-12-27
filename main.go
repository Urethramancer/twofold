package main

import (
	"path/filepath"

	"github.com/grimdork/opt"
)

var opts struct {
	opt.DefaultHelp
	Version  bool   `short:"V" long:"version" description:"Prints version and exits."`
	Verbose  bool   `short:"v" description:"Verbose output. Shows progress of checksumming for each file, not just the list of duplicates."`
	List     bool   `short:"l" long:"list" description:"List duplicates only."`
	Symlink  bool   `long:"symlink" description:"Symlink all duplicates to the first file."`
	Hardlink bool   `long:"hardlink" description:"Hardlink all duplicates to the first file."`
	Remove   bool   `long:"remove" description:"Remove duplicates of the first file."`
	Tens     bool   `short:"s" long:"si" description:"Use SI numbers." default-value:"true"`
	Path     string `placeholder:"DIRECTORY" help:"Directory to look deeply into for duplicates." default-value:"."`
}

func main() {
	a := opt.Parse(&opts)
	if opts.Help {
		a.Usage()
		return
	}

	if opts.Version {
		printVersion()
		return
	}

	dir, err := filepath.Abs(opts.Path)
	if err != nil {
		pr("Error: %s", err.Error())
		return
	}

	if !(opts.List || opts.Hardlink || opts.Symlink || opts.Remove) {
		pr("No flags specified. Bailing out.")
		return
	}

	list := scanDir(dir)
	if opts.Verbose {
		pr("\n")
	}

	if len(list) == 0 {
		pr("No duplicates found.")
		return
	}

	// Only display duplicates
	if opts.List {
		listDuplicates(list)
		return
	}

	// Do something about them
	handleDuplicates(list)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/grimdork/climate/arg"
)

var (
	// mirror the flags in a simple struct for program logic
	flags = struct {
		Verbose  bool
		List     bool
		Symlink  bool
		Hardlink bool
		Remove   bool
		Apply    bool
		Tens     bool
		Path     string
		Workers  int
	}{}
)

func main() {
	opt := arg.New("twofold")
	opt.SetDefaultHelp(true)
	opt.SetOption(arg.GroupDefault, "v", "verbose", "Verbose output. Shows progress of checksumming for each file, not just the list of duplicates.", false, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "l", "list", "List duplicates only.", false, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "", "symlink", "Symlink all duplicates to the first file.", false, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "", "hardlink", "Hardlink all duplicates to the first file.", false, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "", "remove", "Remove duplicates of the first file.", false, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "", "apply", "Apply changes. Without this flag the program runs as a dry-run.", false, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "s", "si", "Use SI numbers.", true, false, arg.VarBool, nil)
	opt.SetOption(arg.GroupDefault, "", "path", "Directory to look deeply into for duplicates.", ".", false, arg.VarString, nil)
	opt.SetOption(arg.GroupDefault, "", "workers", "Number of concurrent hashing workers (0 or 1 = single-threaded). Default: number of CPU cores.", 0, false, arg.VarInt, nil)

	err := opt.Parse(os.Args[1:])
	if err != nil {
		if err == arg.ErrNoArgs {
			opt.PrintHelp()
			return
		}
		fmt.Fprintf(os.Stderr, "Error parsing args: %v\n", err)
		os.Exit(2)
	}

	// read values back
	flags.Verbose = opt.GetBool("verbose")
	flags.List = opt.GetBool("list")
	flags.Symlink = opt.GetBool("symlink")
	flags.Hardlink = opt.GetBool("hardlink")
	flags.Remove = opt.GetBool("remove")
	flags.Apply = opt.GetBool("apply")
	flags.Tens = opt.GetBool("si")
	flags.Path = opt.GetString("path")
	workers := opt.GetInt("workers")
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	flags.Workers = workers

	// normalize path
	dir := flags.Path
	if strings.TrimSpace(dir) == "" {
		dir = "."
	}
	dir, err = filepath.Abs(dir)
	if err != nil {
		pr("Error: %s", err.Error())
		return
	}

	if !(flags.List || flags.Hardlink || flags.Symlink || flags.Remove) {
		pr("No flags specified. Bailing out.")
		return
	}

	list, err := scanDir(dir)
	if err != nil {
		pr("Scan error: %s", err.Error())
		return
	}

	if flags.Verbose {
		pr("\n")
	}

	if len(list) == 0 {
		pr("No duplicates found.")
		return
	}

	// Only display duplicates
	if flags.List {
		listDuplicates(list)
		return
	}

	// Do something about them
	handleDuplicates(list)
}

package main

import "os"

// Duplicates hold the common attributes of matching files, and the list of said files.
type Duplicates struct {
	Hash  string
	Size  int64
	Files *[]string
}

func listDuplicates(list []*Duplicates) {
	for _, d := range list {
		pr("Duplicates of %s (%s):", d.Hash, humanNumber(d.Size, opts.Tens))
		for i, n := range *d.Files {
			pr("\t%d\t%s", i+1, n)
		}
		pr("")
	}
}

func handleDuplicates(list []*Duplicates) {
	for _, d := range list {
		files := *d.Files
		orig := files[0]
		files = files[1:]
		for _, f := range files {
			err := os.Remove(f)
			if err != nil {
				pr("Error: %s", err.Error())
			}
			if opts.Hardlink {
				err = os.Link(orig, f)
				if err != nil {
					pr("Error: %s", err.Error())
				}
				pr("Hardlinked %s to %s", f, orig)
				continue
			}
			if opts.Symlink {
				err = os.Symlink(orig, f)
				if err != nil {
					pr("Error: %s", err.Error())
				}
				pr("Symlinked %s to %s", f, orig)
				continue
			}
			if opts.Remove {
				pr("Removed %s", f)
			}
		}
	}
}

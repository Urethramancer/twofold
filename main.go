package main

import (
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	Verbose  bool `short:"v" description:"Verbose output. Shows progress of checksumming for each file, not just the list of duplicates."`
	List     bool `short:"l" long:"list" description:"List duplicates only."`
	Symlink  bool `long:"symlink" description:"Symlink all duplicates to the first file."`
	Hardlink bool `long:"hardlink" description:"Hardlink all duplicates to the first file."`
	Tens     bool `short:"s" long:"si" description:"Use SI numbers." default-value:"true"`
	Args     struct {
		Path string `required:"true" positional-arg-name:"DIRECTORY" decription:"Directory to look deeply into for duplicates."`
	} `positional-args:"true"`
}

type Duplicates struct {
	Hash  string
	Size  int64
	Files *[]string
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	dir, err := filepath.Abs(opts.Args.Path)
	if err != nil {
		pr("Error: %s", err.Error())
		return
	}

	if opts.List || opts.Hardlink || opts.Symlink {
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
			for _, d := range list {
				pr("Duplicates of %s (%s):", d.Hash, humanNumber(d.Size, opts.Tens))
				for i, n := range *d.Files {
					pr("\t%d\t%s", i+1, n)
				}
				pr("")
			}
			return
		}

		// Do something about them
		for _, d := range list {
			files := *d.Files
			orig := files[0]
			files = files[1:]
			for _, f := range files {
				os.Remove(f)
				if opts.Hardlink {
					err = os.Link(orig, f)
					if err != nil {
						pr("Error: %s", err.Error())
					}
					pr("Hardlinked %s to %s", f, orig)
				}
				if opts.Symlink {
					err = os.Symlink(orig, f)
					if err != nil {
						pr("Error: %s", err.Error())
					}
					pr("Symlinked %s to %s", f, orig)}
				}
			}
		}
	}
}

func scanDir(dir string) []*Duplicates {
	pr("Deep-scanning %s\n", dir)
	dup := make(map[string]*Duplicates)
	err := filepath.Walk(dir, func(path string, fi os.FileInfo, err error) error {
		if !fi.IsDir() {
			hash := hashFile(path, fi.Size())
			if hash == "" {
				return nil
			}

			_, ok := dup[hash]
			if !ok {
				list := make([]string, 0)
				list = append(list, path)
				dup[hash] = &Duplicates{Hash: hash, Size: fi.Size(), Files: &list}
			} else {
				if fi.Size() == dup[hash].Size {
					list := *dup[hash].Files
					list = append(list, path)
					dup[hash].Files = &list
				}
			}
		}
		return nil
	})

	if err != nil {
		pr("Error: %s", err.Error())
		os.Exit(2)
	}

	list := make([]*Duplicates, 0)
	for _, v := range dup {
		if len(*v.Files) > 1 {
			list = append(list, v)
		}
	}

	return list
}

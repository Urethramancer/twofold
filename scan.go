package main

import (
	"os"
	"path/filepath"
)

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
				var list []string
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

	var list []*Duplicates
	for _, v := range dup {
		if len(*v.Files) > 1 {
			list = append(list, v)
		}
	}

	return list
}

package main

import (
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

func scanDir(dir string) ([]*Duplicates, error) {
	pr("Deep-scanning %s\n", dir)

	// If workers <= 1, fall back to simple synchronous path
	if flags.Workers <= 1 {
		dup := make(map[string]*Duplicates)
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				pr("walk error: %s", err.Error())
				return nil
			}

			// skip .git directories and symlinked directories
			if d.IsDir() {
				if d.Name() == ".git" {
					return filepath.SkipDir
				}
				if d.Type()&fs.ModeSymlink != 0 {
					return filepath.SkipDir
				}
				return nil
			}

			fi, err := d.Info()
			if err != nil {
				pr("failed to stat %s: %s", path, err.Error())
				return nil
			}

			hash := hashFile(path, fi.Size())
			if hash == "" {
				return nil
			}

			entry, ok := dup[hash]
			if !ok {
				dup[hash] = &Duplicates{Hash: hash, Size: fi.Size(), Files: []string{path}}
			} else {
				if fi.Size() == entry.Size {
					entry.Files = append(entry.Files, path)
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

		var list []*Duplicates
		for _, v := range dup {
			if len(v.Files) > 1 {
				list = append(list, v)
			}
		}

		return list, nil
	}

	// Concurrent path with aggregated progress
	// First, count files to hash so the progress bar has a total
	var total int64 = 0
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			if d.Type()&fs.ModeSymlink != 0 {
				return filepath.SkipDir
			}
			return nil
		}
		total++
		return nil
	})

	bar := pb.New64(total)
	if flags.Verbose {
		bar.Set(pb.Bytes, false)
		bar.Start()
	}

	type job struct {
		path string
		size int64
	}
	type result struct {
		path string
		size int64
		hash string
	}

	jobs := make(chan job, flags.Workers*2)
	results := make(chan result, flags.Workers*2)

	var wg sync.WaitGroup
	// start workers
	for i := 0; i < flags.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				h := hashFile(j.path, j.size)
				if h == "" {
					if flags.Verbose {
						bar.Add(1)
					}
					continue
				}
				results <- result{path: j.path, size: j.size, hash: h}
				if flags.Verbose {
					bar.Add(1)
				}
			}
		}()
	}

	// producer: walk directory and feed jobs
	go func() {
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				pr("walk error: %s", err.Error())
				return nil
			}
			if d.IsDir() {
				if d.Name() == ".git" {
					return filepath.SkipDir
				}
				if d.Type()&fs.ModeSymlink != 0 {
					return filepath.SkipDir
				}
				return nil
			}
			fi, err := d.Info()
			if err != nil {
				pr("failed to stat %s: %s", path, err.Error())
				return nil
			}
			jobs <- job{path: path, size: fi.Size()}
			return nil
		})
		close(jobs)
	}()

	// collector
	go func() {
		wg.Wait()
		close(results)
		if flags.Verbose {
			bar.Finish()
		}
	}()

	dup := make(map[string]*Duplicates)
	for r := range results {
		entry, ok := dup[r.hash]
		if !ok {
			dup[r.hash] = &Duplicates{Hash: r.hash, Size: r.size, Files: []string{r.path}}
		} else {
			if r.size == entry.Size {
				entry.Files = append(entry.Files, r.path)
			}
		}
	}

	var list []*Duplicates
	for _, v := range dup {
		if len(v.Files) > 1 {
			list = append(list, v)
		}
	}

	return list, nil
}

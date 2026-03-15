package main

import (
	"io"
	"os"
	"strings"
)

// linkFunc is a variable wrapper around os.Link so tests can override link behaviour.
var linkFunc = os.Link

// Duplicates hold the common attributes of matching files, and the list of said files.
type Duplicates struct {
	Hash  string
	Size  int64
	Files []string
}

func listDuplicates(list []*Duplicates) {
	for _, d := range list {
		pr("Duplicates of %s (%s):", d.Hash, humanNumber(d.Size, flags.Tens))
		for i, n := range d.Files {
			pr("\t%d\t%s", i+1, n)
		}
		pr("")
	}
}

func handleDuplicates(list []*Duplicates) {
	for _, d := range list {
		files := d.Files
		if len(files) < 2 {
			continue
		}
		orig := files[0]
		duplicates := files[1:]
		for _, f := range duplicates {
			// dry-run by default: only show what we'd do
			if !flags.Apply {
				if flags.Hardlink {
					pr("Would hardlink %s -> %s", f, orig)
					continue
				}
				if flags.Symlink {
					pr("Would symlink %s -> %s", f, orig)
					continue
				}
				if flags.Remove {
					pr("Would remove %s", f)
					continue
				}
				// Shouldn't reach here; no action requested
				continue
			}

			// Apply mode: perform safe replacement
			if flags.Hardlink {
				// attempt to create hardlink at destination pointing to orig
				_ = os.Remove(f)
				err := linkFunc(orig, f)
				if err != nil {
					// fallback: cross-device link not possible — copy the file instead
					if isCrossDevice(err) {
						cpErr := copyFile(orig, f)
						if cpErr != nil {
							pr("Error copying across devices %s -> %s: %s", orig, f, cpErr.Error())
							continue
						}
						pr("Copied %s to %s (cross-device)", orig, f)
						continue
					}
					pr("Error hardlinking %s -> %s: %s", f, orig, err.Error())
					continue
				}
				pr("Hardlinked %s to %s", f, orig)
				continue
			}

			if flags.Symlink {
				// create symlink at path f pointing to orig
				err := os.Remove(f)
				if err != nil && !os.IsNotExist(err) {
					pr("Error removing existing %s before symlink: %s", f, err.Error())
					continue
				}
				err = os.Symlink(orig, f)
				if err != nil {
					pr("Error symlinking %s -> %s: %s", f, orig, err.Error())
					continue
				}
				pr("Symlinked %s to %s", f, orig)
				continue
			}

			if flags.Remove {
				err := os.Remove(f)
				if err != nil {
					pr("Error removing %s: %s", f, err.Error())
					continue
				}
				pr("Removed %s", f)
				continue
			}
		}
	}
}

// copyFile copies file contents and preserves mode
func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	if err != nil {
		return err
	}

	// preserve mode
	st, err := os.Stat(src)
	if err == nil {
		_ = os.Chmod(dst, st.Mode())
	}
	return nil
}

// isCrossDevice returns true for EXDEV link errors
func isCrossDevice(err error) bool {
	if err == nil {
		return false
	}
	typ := err.Error()
	return strings.Contains(typ, "invalid cross-device link") || strings.Contains(typ, "cross-device link") || strings.Contains(typ, "EXDEV")
}

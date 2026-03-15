package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// helper to create files with same content
func createFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func TestScanDir_SyncAndConcurrent(t *testing.T) {
	dir := t.TempDir()
	// create three files where two are duplicates
	f1 := createFile(t, dir, "a.txt", "hello world")
	f2 := createFile(t, dir, "b.txt", "hello world")
	createFile(t, dir, "c.txt", "other")

	// synchronous
	flags.Workers = 1
	list, err := scanDir(dir)
	if err != nil {
		t.Fatalf("scanDir sync failed: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 duplicate group, got %d", len(list))
	}
	grp := list[0]
	if grp.Size == 0 {
		t.Fatalf("group size zero")
	}
	if len(grp.Files) != 2 {
		t.Fatalf("expected 2 files in group, got %d", len(grp.Files))
	}

	// concurrent: use number of CPUs (at least 2 if available)
	workers := runtime.NumCPU()
	if workers < 2 {
		workers = 2
	}
	flags.Workers = workers
	list2, err := scanDir(dir)
	if err != nil {
		t.Fatalf("scanDir concurrent failed: %v", err)
	}
	if len(list2) != 1 {
		t.Fatalf("expected 1 duplicate group (concurrent), got %d", len(list2))
	}
	// check that file paths are present
	found := map[string]bool{}
	for _, p := range list2[0].Files {
		found[p] = true
	}
	if !found[f1] || !found[f2] {
		t.Fatalf("expected both duplicate paths present")
	}
}

func TestHandleDuplicates_RemoveAndSymlink(t *testing.T) {
	dir := t.TempDir()
	orig := createFile(t, dir, "orig.txt", "same")
	dup := createFile(t, dir, "dup.txt", "same")

	grp := &Duplicates{Hash: "h", Size: 4, Files: []string{orig, dup}}

	// dry-run should not remove
	flags.Apply = false
	flags.Remove = true
	handleDuplicates([]*Duplicates{grp})
	if _, err := os.Stat(dup); os.IsNotExist(err) {
		t.Fatalf("dup removed during dry-run")
	}

	// Apply remove
	flags.Apply = true
	flags.Remove = true
	handleDuplicates([]*Duplicates{grp})
	if _, err := os.Stat(dup); !os.IsNotExist(err) {
		t.Fatalf("dup not removed in apply mode: %v", err)
	}

	// recreate duplicate to test symlink
	dup2 := createFile(t, dir, "dup2.txt", "same")
	grp2 := &Duplicates{Hash: "h2", Size: 4, Files: []string{orig, dup2}}
	flags.Apply = true
	flags.Remove = false
	flags.Symlink = true
	handleDuplicates([]*Duplicates{grp2})
	// dup2 should now be a symlink
	st, err := os.Lstat(dup2)
	if err != nil {
		t.Fatalf("stat dup2: %v", err)
	}
	if st.Mode()&os.ModeSymlink == 0 {
		t.Fatalf("dup2 is not a symlink")
	}
}

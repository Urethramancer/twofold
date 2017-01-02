package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

func hashFile(name string, size int64) string {
	file, err := os.Open(name)
	if err != nil {
		pr("Couldn't open %s", name)
		return ""
	}

	src := io.Reader(file)
	if opts.Verbose {
		pr("%s", name)
		bar := pb.StartNew(int(size))
		bar.ShowSpeed = true
		bar.SetUnits(pb.U_BYTES)
		bar.ShowTimeLeft = false
		src = bar.NewProxyReader(file)
		defer bar.Update()
	}

	hash := sha256.New()
	_, err = io.Copy(hash, src)
	file.Close()
	if err != nil {
		pr("Couldn't read/hash %s", name)
		return ""
	}

	return hex.EncodeToString(hash.Sum(nil))
}

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

func hashFile(name string, size int64) string {
	file, err := os.Open(name)
	if err != nil {
		pr("Couldn't open %s", name)
		return ""
	}
	defer file.Close()

	src := io.Reader(file)
	if flags.Verbose {
		pr("%s", name)
		// pb/v3 uses NewPool or default bar; use simple default bar
		bar := pb.New64(0)
		bar.Set(pb.Bytes, true)
		bar.Start()
		src = bar.NewProxyReader(file)
		defer bar.Finish()
	}

	h := sha256.New()
	_, err = io.Copy(h, src)
	if err != nil {
		pr("Couldn't read/hash %s", name)
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

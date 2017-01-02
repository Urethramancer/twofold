package main

import (
	"fmt"
	"math"
	"os"
)

func pr(f string, v ...interface{}) {
	fmt.Printf(f+"\n", v...)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func humanNumber(n int64, tens bool) string {
	num := float64(n)
	var unit float64 = 1024
	if tens {
		unit = 1000
	}
	if num < unit {
		return fmt.Sprintf("%dB", int(num))
	}
	exp := int(math.Log(num) / math.Log(unit))
	pre := "kMGTPE"
	pre = pre[exp-1 : exp]
	if !tens {
		pre = pre + "i"
	}
	r := uint64(n) / uint64(math.Pow(unit, float64(exp)))
	return fmt.Sprintf("%d %sB", r, pre)
}

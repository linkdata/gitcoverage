package main

import (
	"flag"
	"testing"
)

func TestMainFn(t *testing.T) {
	flag.Parse()
	mainfn()
}

func TestMainError(t *testing.T) {
	exitFn = func(i int) {
		if i == 0 {
			t.Error(i)
		}
		if i == 125 {
			t.Log("didn't get a syscall.Errno")
		}
	}
	flag.Parse()
	*flagDebug = true
	main()
}

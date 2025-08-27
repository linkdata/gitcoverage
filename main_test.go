package main

import (
	"flag"
	"os"
	"testing"
)

func TestMainFn(t *testing.T) {
	flag.Parse()
	mainfn()
}

func TestMainError(t *testing.T) {
	os.Args = append(os.Args, "/proc/.nonexistant")
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

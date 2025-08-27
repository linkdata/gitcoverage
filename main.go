package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"

	gitsemver "github.com/linkdata/gitsemver/pkg"
)

var (
	flagGit   = flag.String("git", "git", "path to Git executable")
	flagDebug = flag.Bool("debug", false, "write debug info to stderr")
	flagPct   = flag.String("pct", "", "coverage percentage (required)")
)

func mainfn() int {
	var err error
	if *flagPct != "" {
		repoDir := os.ExpandEnv(flag.Arg(0))
		if repoDir == "" {
			repoDir = "."
		}
		var vs *gitsemver.GitSemVer
		if vs, err = gitsemver.New(*flagGit); err == nil {
			if *flagDebug {
				vs.DebugOut = os.Stderr
			}
			if repoDir, err = vs.Git.CheckGitRepo(repoDir); err == nil {
				var vi gitsemver.VersionInfo
				if vi, err = vs.GetVersion(repoDir); err == nil {

					fmt.Println(vi.Branch)
					return 0
				}
			}
		}
	}

	retv := 125
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if e := errors.Unwrap(err); e != nil {
			if errno, ok := e.(syscall.Errno); ok {
				retv = int(errno)
			}
		}
	}
	return retv
}

var exitFn func(int) = os.Exit

func main() {
	flag.Parse()
	exitFn(mainfn())
}

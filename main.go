package main

import (
	"os"

	"github.com/lingrino/vaku/v2/cmd"
)

// version is populated at build time by goreleaser.
var version = "dev"

// used for testing.
var executeCMD = cmd.Execute
var exitCmd = os.Exit

func main() {
	code := executeCMD(version, os.Args[1:], os.Stdout, os.Stderr)
	exitCmd(code)
}

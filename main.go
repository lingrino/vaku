package main

import (
	"os"

	"github.com/lingrino/vaku/cmd"
)

// version is populated at build time by goreleaser
var version = "dev"

var executeCMD = cmd.Execute
var exitCmd = os.Exit

func main() {
	exitCmd(executeCMD(version))
}

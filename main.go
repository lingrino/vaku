package main

import (
	"github.com/lingrino/vaku/cmd"
)

// version is populated at build time by goreleaser
var version = "dev"

var executeCMD = cmd.Execute

func main() {
	executeCMD(version)
}

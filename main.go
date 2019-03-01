package main

import (
	"github.com/lingrino/vaku/cmd"
)

// version is populated at build time by goreleaser
var version = "dev"

func main() {
	cmd.Execute(version)
}

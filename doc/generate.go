package main

import (
	"log"

	"github.com/Lingrino/vaku/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.VakuCmd, "./doc/")
	if err != nil {
		log.Fatal(err)
	}
}

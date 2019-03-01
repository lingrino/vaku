package main

import (
	"log"

	"github.com/lingrino/vaku/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.VakuCmd, "./docs/")
	if err != nil {
		log.Fatal(err)
	}
}

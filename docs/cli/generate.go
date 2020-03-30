package main

import (
	"fmt"
	"os"

	"github.com/lingrino/vaku/cmd"
	"github.com/spf13/cobra/doc"
)

const (
	exitFail = 1
)

// main calls generate.
// from repo root - go run docs/cli
func main() {
	err := generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

// generate generates documentation for the cli into the docs folder.
// https://github.com/spf13/cobra/blob/master/doc/md_docs.md
func generate() error {
	vc, err := cmd.NewVakuCmd("version")
	if err != nil {
		return err
	}

	err = doc.GenMarkdownTree(vc, "./docs/cli/")
	if err != nil {
		return err
	}

	return nil
}

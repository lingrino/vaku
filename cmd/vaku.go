package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	exitFail = 1
)

const (
	vakuUse     = "vaku <cmd>"
	vakuShort   = "Vaku is a tool for working with large vault k/v secret engines"
	vakuExample = "vaku folder list secret/foo"
	vakuLong    = `vaku
long
description

CLI documentation - 'vaku help [cmd]'
API documentation - https://pkg.go.dev/github.com/lingrino/vaku/vaku
Built by Sean Lingren <sean@lingrino.com>`
)

func newVakuCmd(version string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     vakuUse,
		Short:   vakuShort,
		Long:    vakuLong,
		Example: vakuExample,

		// https://github.com/spf13/cobra/issues/914#issuecomment-548411337
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.AddCommand(
		newCompletionCmd(),
		newDocsCmd(),
		newPathCmd(),
		newFolderCmd(),
		newVersionCmd(version),
	)

	return cmd, nil
}

// Execute runs Vaku
func Execute(version string) {
	vc, err := newVakuCmd(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(exitFail)
	}

	err = vc.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(exitFail)
	}
}

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
	vakuUse   = "vaku"
	vakuShort = "Vaku is a tool for working with large vault k/v secret engines"
	vakuLong  = `vaku
long
description

CLI documentation - 'vaku help [cmd]'
API documentation - https://pkg.go.dev/github.com/lingrino/vaku/vaku
Built by Sean Lingren <sean@lingrino.com>`
)

func NewVakuCmd(version string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   vakuUse,
		Short: vakuShort,
		Long:  vakuLong,
	}

	cmd.AddCommand(
		newPathCmd(),
		newFolderCmd(),
		newVersionCmd(version),
	)

	return cmd, nil
}

// Execute runs Vaku
func Execute(version string) {
	vc, err := NewVakuCmd(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}

	err = vc.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

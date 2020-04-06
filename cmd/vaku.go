package cmd

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
)

const (
	exitSuccess = 0
	exitFail    = 1
)

// Special strings used for failure injection in tests
const (
	failString = "fail"
	testString = "test"
)

const (
	vakuUse     = "vaku <cmd>"
	vakuShort   = "Vaku is a CLI for working with large Vault k/v secret engines"
	vakuExample = "vaku folder list secret/foo"
	vakuLong    = `Vaku is a CLI for working with large Vault k/v secret engines

The Vaku CLI provides path and folder based commands that work on
both version 1 and version 2 k/v secret engines. Vaku can help manage
large amounts of vault data by updating secrets in place, moving
paths or folders, searching secrets, and more.

Vaku is not a replacement for the Vault CLI and requires that you
already are authenticated to Vault before running any commands. Vaku
commands should not be run on non-k/v engines.

CLI documentation - 'vaku help [cmd]'
API documentation - https://pkg.go.dev/github.com/lingrino/vaku/vaku
Built by Sean Lingren <sean@lingrino.com>`
)

func newVakuCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     vakuUse,
		Short:   vakuShort,
		Long:    vakuLong,
		Example: vakuExample,

		// https://github.com/spf13/cobra/issues/914#issuecomment-548411337
		SilenceErrors: true,
		SilenceUsage:  true,

		DisableAutoGenTag: true,
	}

	cmd.AddCommand(
		newCompletionCmd(),
		newDocsCmd(),
		newPathCmd(),
		newFolderCmd(),
		newVersionCmd(version),
	)

	return cmd
}

// Execute runs Vaku
func Execute(version string) int {
	vc := newVakuCmd(version)

	// Test/Failure injection
	if version == testString || version == failString {
		var nilout bytes.Buffer
		vc.SetOut(&nilout)
		vc.SetErr(&nilout)
	}
	if version == failString {
		vc.SetArgs([]string{failString})
	}

	err := vc.Execute()
	if err != nil {
		fmt.Fprintf(vc.ErrOrStderr(), "Error: %s\n", err)
		return exitFail
	}

	return exitSuccess
}

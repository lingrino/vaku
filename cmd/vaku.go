package cmd

import (
	"github.com/spf13/cobra"
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

// newVakuCmd sets flags/subcommands and returns the base vaku command.
func (c *cli) newVakuCmd() *cobra.Command {
	// base command
	cmd := &cobra.Command{
		Use:     vakuUse,
		Short:   vakuShort,
		Long:    vakuLong,
		Example: vakuExample,

		PersistentPreRunE: c.validateVakuFlags,

		// https://github.com/spf13/cobra/issues/914#issuecomment-548411337
		SilenceErrors: true,
		SilenceUsage:  true,

		// prevents docs from adding promotional message footer
		DisableAutoGenTag: true,
	}

	// add base/persistent flags
	c.addVakuFlags(cmd)

	// add subcommands
	cmd.AddCommand(
		c.newCompletionCmd(),
		c.newDocsCmd(),
		c.newFolderCmd(),
		c.newPathCmd(),
		c.newVersionCmd(),
	)

	return cmd
}

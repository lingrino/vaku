package cmd

import (
	"github.com/spf13/cobra"
)

const (
	vakuUse     = "vaku <cmd>"
	vakuShort   = "Vaku is a CLI for working with large Vault k/v secret engines"
	vakuExample = "vaku folder list secret/foo"
	vakuLong    = `Vaku is a CLI for working with large Vault k/v secret engines

The Vaku CLI provides path- and folder-based commands that work on
both Version 1 and Version 2 K/V secret engines. Vaku can help manage
large amounts of Vault data by updating secrets in place, moving
paths or folders, searching secrets, and more.

Vaku is not a replacement for the Vault CLI and requires that you are
already authenticated to Vault before running any commands. Vaku
commands should not be run on non-K/V engines.

CLI documentation - 'vaku help [cmd]'
API documentation - https://pkg.go.dev/github.com/lingrino/vaku/v2/api
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
		c.newDocsCmd(),
		c.newFolderCmd(),
		c.newPathCmd(),
		c.newVersionCmd(),
	)

	return cmd
}

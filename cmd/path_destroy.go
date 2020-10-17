package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathDestroyUse   = "destroy"
	pathDestroyShort = "Vaku CLI does not yet support path destroy. Use the vaku API or native Vault CLI"
	pathDestroyLong  = "Vaku CLI does not yet support path destroy. Use the vaku API or native Vault CLI"
)

func (c *cli) newPathDestroyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pathDestroyUse,
		Short: pathDestroyShort,
		Long:  pathDestroyLong,

		// disable all discovery
		Hidden:                true,
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		PersistentPreRun:      nil,
		Args:                  cobra.ArbitraryArgs,
	}

	return cmd
}

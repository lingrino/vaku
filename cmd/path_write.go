package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathWriteUse   = "write"
	pathWriteShort = "Vaku CLI does not support path write. Use the vaku API or native Vault CLI"
	pathWriteLong  = "Vaku CLI does not support path write. Use the vaku API or native Vault CLI"
)

func (c *cli) newPathWriteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pathWriteUse,
		Short: pathWriteShort,
		Long:  pathWriteLong,

		// disable all discovery
		Hidden:                true,
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		PersistentPreRun:      nil,
		Args:                  cobra.ArbitraryArgs,
	}

	return cmd
}

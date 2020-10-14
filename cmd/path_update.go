package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathUpdateUse   = "update"
	pathUpdateShort = "Vaku CLI does not yet support path update. Use the vaku API or native Vault CLI"
	pathUpdateLong  = "Vaku CLI does not yet support path update. Use the vaku API or native Vault CLI"
)

func (c *cli) newPathUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   pathUpdateUse,
		Short: pathUpdateShort,
		Long:  pathUpdateLong,

		// disable all discovery
		Hidden:                true,
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		PersistentPreRun:      nil,
		Args:                  cobra.ArbitraryArgs,
	}

	return cmd
}

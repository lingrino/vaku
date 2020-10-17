package cmd

import (
	"github.com/spf13/cobra"
)

const (
	folderDestroyUse   = "destroy"
	folderDestroyShort = "Vaku CLI does not yet support folder destroy. Use the vaku API or native Vault CLI"
	folderDestroyLong  = "Vaku CLI does not yet support folder destroy. Use the vaku API or native Vault CLI"
)

func (c *cli) newFolderDestroyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   folderDestroyUse,
		Short: folderDestroyShort,
		Long:  folderDestroyLong,

		// disable all discovery
		Hidden:                true,
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		PersistentPreRun:      nil,
		Args:                  cobra.ArbitraryArgs,
	}

	return cmd
}

package cmd

import (
	"github.com/spf13/cobra"
)

const (
	folderWriteUse   = "write"
	folderWriteShort = "Vaku CLI does not support folder write. Use the vaku API or native Vault CLI"
	folderWriteLong  = "Vaku CLI does not support folder write. Use the vaku API or native Vault CLI"
)

func (c *cli) newFolderWriteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   folderWriteUse,
		Short: folderWriteShort,
		Long:  folderWriteLong,

		// disable all discovery
		Hidden:                true,
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		PersistentPreRun:      nil,
		Args:                  cobra.ArbitraryArgs,
	}

	return cmd
}

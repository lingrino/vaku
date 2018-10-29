package cmd

import (
	"github.com/spf13/cobra"
)

var folderCmd = &cobra.Command{
	Use:   "folder [cmd]",
	Short: "Contains all vaku folder functions, does nothing on its own",

	// Auth to vault on all commands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		authVGC()
	},
}

func init() {
	VakuCmd.AddCommand(folderCmd)
}

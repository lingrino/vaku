package cmd

import (
	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path [cmd]",
	Short: "Contains all vaku path functions, does nothing on its own",

	// Auth to vault on all commands
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		authVGC()
	},
}

func init() {
	VakuCmd.AddCommand(pathCmd)
}

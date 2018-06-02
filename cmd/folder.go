package cmd

import (
	"github.com/spf13/cobra"
)

var folderCmd = &cobra.Command{
	Use:   "folder [cmd]",
	Short: "Contains the vaku folder functions, does nothing on its own",
}

func init() {
	rootCmd.AddCommand(folderCmd)
}

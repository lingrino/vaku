package cmd

import (
	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path [cmd]",
	Short: "Contains the vaku path functions, does nothing on its own",
}

func init() {
	vakuCmd.AddCommand(pathCmd)
}

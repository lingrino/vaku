package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the current Vaku CLI and API versions",

	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CLI Version: 1.0")
		fmt.Println("API Version:", vaku.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

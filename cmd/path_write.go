package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pathWriteCmd = &cobra.Command{
	Hidden:                true,
	DisableSuggestions:    true,
	DisableFlagsInUseLine: true,

	Use:   "write",
	Short: "Vaku CLI does not support writes. Please use either the native Vault CLI or the Vaku API",
	Long:  "Vaku CLI does not support writes. Please use either the native Vault CLI or the Vaku API",

	Args:             cobra.ArbitraryArgs,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("ERROR: Vaku CLI does not support writes. Please use either the native Vault CLI or the Vaku API")
	},
}

func init() {
	pathCmd.AddCommand(pathWriteCmd)
}

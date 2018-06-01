package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaku",
	Short: "Vaku CLI extends the official Vault CLI with useful high-level functions",
	Long: `Vaku CLI extends the official Vault CLI with useful high-level functions

Built by Sean Lingren <srlingren@gmail.com>
CLI documentation is available using --help
API documentation is available at https://godoc.org/github.com/Lingrino/vaku/vaku`,
}

// Execute initializes and runs the vaku command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

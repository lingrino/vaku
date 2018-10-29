package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

var VakuCmd = &cobra.Command{
	Use:   "vaku",
	Short: "Vaku CLI extends the official Vault CLI with useful high-level functions",
	Long: `Vaku CLI extends the official Vault CLI with useful high-level functions

The Vaku CLI is intended to be used side by side with the official Vault CLI,
and only provides functions to extend the existing functionality. Many of the 'vaku path'
functions are very similar (or even less featured) than the vault CLI equivalent. However
vaku works on both v1 and v2 secret mounts, and can even copy/move secrets between them.

Vaku does not log you in to vault or help you with getting a token. Like the CLI,
it will look for a token first at the VAULT_TOKEN env var and then in ~/.vault-token

Built by Sean Lingren <srlingren@gmail.com>
CLI documentation is available using 'vaku help [cmd]'
API documentation is available at https://godoc.org/github.com/Lingrino/vaku/vaku`,
}

func init() {
	VakuCmd.PersistentFlags().StringVarP(&format, "format", "o", "json", "The output format to use. One of: \"json\", \"text\"")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the VakuCmd.
func Execute(v string) {
	var err error

	version = v

	err = VakuCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

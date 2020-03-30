package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	exitFail = 1
)

var version = "dev"

var VakuCmd = &cobra.Command{
	Use:   "vaku",
	Short: "short description",
	Long: `long description

long stuff

CLI documentation - 'vaku help [cmd]'
API documentation - https://pkg.go.dev/github.com/lingrino/vaku/vaku
Built by Sean Lingren <sean@lingrino.com>`,
}

// Execute runs Vaku
func Execute(v string) {
	version = v

	err := VakuCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitFail)
	}
}

package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathUse     = "path <cmd>"
	pathShort   = "Commands that act on Vault paths"
	pathExample = "vaku path list secret/foo"
	pathLong    = `Commands that act on Vault paths

Commands under the path subcommand act on Vault paths. Vaku can list,
copy, move, search, etc.. on Vault paths.`
)

func newPathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathUse,
		Short:   pathShort,
		Long:    pathLong,
		Example: pathExample,
	}

	return cmd
}

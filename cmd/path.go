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

func (c *cli) newPathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathUse,
		Short:   pathShort,
		Long:    pathLong,
		Example: pathExample,

		PersistentPreRunE: c.initVakuClient,
	}

	c.addPathFolderFlags(cmd)

	cmd.AddCommand(
		c.newPathListCmd(),
		c.newPathReadCmd(),
		c.newPathWriteCmd(),
		c.newPathDeleteCmd(),
		c.newPathSearchCmd(),
	)

	return cmd
}

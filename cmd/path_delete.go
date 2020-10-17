package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathDeleteArgs    = 1
	pathDeleteUse     = "delete <path>"
	pathDeleteShort   = "Delete a secret at a path"
	pathDeleteLong    = "Delete a secret at a path"
	pathDeleteExample = "vaku path delete secret/foo"
)

func (c *cli) newPathDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathDeleteUse,
		Short:   pathDeleteShort,
		Long:    pathDeleteLong,
		Example: pathDeleteExample,

		Args: cobra.ExactArgs(pathDeleteArgs),

		RunE: c.runPathDelete,
	}

	return cmd
}

func (c *cli) runPathDelete(cmd *cobra.Command, args []string) error {
	return c.vc.PathDelete(args[0])
}

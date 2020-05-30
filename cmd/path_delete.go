package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathDeleteUse     = "delete <path>"
	pathDeleteShort   = "Delete all paths at a path"
	pathDeleteExample = "vaku path delete secret/foo"
	pathDeleteLong    = "Delete all paths at a path"
)

func (c *cli) newPathDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathDeleteUse,
		Short:   pathDeleteShort,
		Long:    pathDeleteLong,
		Example: pathDeleteExample,

		Args: cobra.ExactArgs(1),

		RunE: c.runPathDelete,
	}

	return cmd
}

func (c *cli) runPathDelete(cmd *cobra.Command, args []string) error {
	err := c.vc.PathDelete(args[0])
	return err
}

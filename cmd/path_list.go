package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathListArgs    = 1
	pathListUse     = "list <path>"
	pathListShort   = "List all paths at a path"
	pathListLong    = "List all paths at a path"
	pathListExample = "vaku path list secret/foo"
)

func (c *cli) newPathListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathListUse,
		Short:   pathListShort,
		Long:    pathListLong,
		Example: pathListExample,

		Args: cobra.ExactArgs(pathListArgs),

		RunE: c.runPathList,
	}

	return cmd
}

func (c *cli) runPathList(cmd *cobra.Command, args []string) error {
	list, err := c.vc.PathList(args[0])
	c.output(list)
	return err
}

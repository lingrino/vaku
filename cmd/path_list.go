package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathListUse     = "list <path>"
	pathListShort   = "List all paths at a path"
	pathListExample = "vaku path list secret/foo"
	pathListLong    = "List all paths at a path"
)

func (c *cli) newPathListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathListUse,
		Short:   pathListShort,
		Long:    pathListLong,
		Example: pathListExample,

		Args: cobra.ExactArgs(1),

		// RunE: c.runPathList,
	}

	return cmd
}

// func (c *cli) runPathList(cmd *cobra.Command, args []string) error {
// 	list, err := c.vc.PathList(args[0])
// 	c.output(list)
// 	return err
// }

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderListUse     = "list <path>"
	folderListShort   = "Recursively list all paths at a path"
	folderListExample = "vaku folder list secret/foo"
	folderListLong    = "Recursively list all paths at a path"
)

func (c *cli) newFolderListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderListUse,
		Short:   folderListShort,
		Long:    folderListLong,
		Example: folderListExample,

		Args: cobra.ExactArgs(1),

		RunE: c.runfolderList,
	}

	return cmd
}

func (c *cli) runfolderList(cmd *cobra.Command, args []string) error {
	list, err := c.vc.FolderList(context.Background(), args[0])
	c.output(list)
	return err
}

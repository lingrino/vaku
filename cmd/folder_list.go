package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderListArgs    = 1
	folderListUse     = "list <folder>"
	folderListShort   = "Recursively list all paths in a folder"
	folderListLong    = "Recursively list all paths in a folder"
	folderListExample = "vaku folder list secret/foo"
)

func (c *cli) newFolderListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderListUse,
		Short:   folderListShort,
		Long:    folderListLong,
		Example: folderListExample,

		Args: cobra.ExactArgs(folderListArgs),

		RunE: c.runfolderList,
	}

	return cmd
}

func (c *cli) runfolderList(cmd *cobra.Command, args []string) error {
	list, err := c.vc.FolderList(context.Background(), args[0])
	c.output(list)
	return err
}

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderDeleteUse     = "delete <path>"
	folderDeleteShort   = "Recursively delete all paths in a folder"
	folderDeleteExample = "vaku folder delete secret/foo"
	folderDeleteLong    = "Recursively delete all paths in a folder"
)

func (c *cli) newFolderDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderDeleteUse,
		Short:   folderDeleteShort,
		Long:    folderDeleteLong,
		Example: folderDeleteExample,

		Args: cobra.ExactArgs(1),

		RunE: c.runfolderDelete,
	}

	return cmd
}

func (c *cli) runfolderDelete(cmd *cobra.Command, args []string) error {
	return c.vc.FolderDelete(context.Background(), args[0])
}

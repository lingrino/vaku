package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderDeleteArgs    = 1
	folderDeleteUse     = "delete <folder>"
	folderDeleteShort   = "Recursively delete all secrets in a folder"
	folderDeleteLong    = "Recursively delete all secrets in a folder"
	folderDeleteExample = "vaku folder delete secret/foo"
)

func (c *cli) newFolderDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderDeleteUse,
		Short:   folderDeleteShort,
		Long:    folderDeleteLong,
		Example: folderDeleteExample,

		Args: cobra.ExactArgs(folderDeleteArgs),

		RunE: c.runfolderDelete,
	}

	return cmd
}

func (c *cli) runfolderDelete(cmd *cobra.Command, args []string) error {
	return c.vc.FolderDelete(context.Background(), args[0])
}

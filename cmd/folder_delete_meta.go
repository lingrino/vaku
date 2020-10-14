package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderDeleteMetaArgs    = 1
	folderDeleteMetaUse     = "delete-meta <folder>"
	folderDeleteMetaShort   = "Recursively delete all secrets metadata and versions in a folder"
	folderDeleteMetaLong    = "Recursively delete all secrets metadata and versions in a folder"
	folderDeleteMetaExample = "vaku folder delete-meta secret/foo"
)

func (c *cli) newFolderDeleteMetaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderDeleteMetaUse,
		Short:   folderDeleteMetaShort,
		Long:    folderDeleteMetaLong,
		Example: folderDeleteMetaExample,

		Args: cobra.ExactArgs(folderDeleteMetaArgs),

		RunE: c.runfolderDeleteMeta,
	}

	return cmd
}

func (c *cli) runfolderDeleteMeta(cmd *cobra.Command, args []string) error {
	return c.vc.FolderDeleteMeta(context.Background(), args[0])
}

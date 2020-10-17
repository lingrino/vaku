package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderCopyArgs    = 2
	folderCopyUse     = "copy <source folder> <destination folder>"
	folderCopyShort   = "Recursively copy all secrets in source folder to destination folder"
	folderCopyLong    = "Recursively copy all secrets in source folder to destination folder"
	folderCopyExample = "vaku folder copy secret/foo secret/bar"
)

func (c *cli) newFolderCopyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderCopyUse,
		Short:   folderCopyShort,
		Long:    folderCopyLong,
		Example: folderCopyExample,

		Args: cobra.ExactArgs(folderCopyArgs),

		RunE: c.runfolderCopy,
	}

	return cmd
}

func (c *cli) runfolderCopy(cmd *cobra.Command, args []string) error {
	return c.vc.FolderCopy(context.Background(), args[0], args[1])
}

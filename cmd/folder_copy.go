package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderCopyUse     = "copy <source folder> <destination folder>"
	folderCopyShort   = "Copy a folder from source to destination"
	folderCopyExample = "vaku folder copy secret/foo secret/bar"
	folderCopyLong    = "Copy a folder from source to destination"
)

func (c *cli) newFolderCopyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderCopyUse,
		Short:   folderCopyShort,
		Long:    folderCopyLong,
		Example: folderCopyExample,

		Args: cobra.ExactArgs(2), //nolint:gomnd

		RunE: c.runfolderCopy,
	}

	return cmd
}

func (c *cli) runfolderCopy(cmd *cobra.Command, args []string) error {
	return c.vc.FolderCopy(context.Background(), args[0], args[1])
}

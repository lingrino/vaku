package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderMoveArgs    = 2
	folderMoveUse     = "move <source folder> <destination folder>"
	folderMoveShort   = "Recursively move all secrets in source folder to destination folder"
	folderMoveLong    = "Recursively move all secrets in source folder to destination folder"
	folderMoveExample = "vaku folder move secret/foo secret/bar"
)

func (c *cli) newFolderMoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderMoveUse,
		Short:   folderMoveShort,
		Long:    folderMoveLong,
		Example: folderMoveExample,

		Args: cobra.ExactArgs(folderMoveArgs),

		RunE: c.runfolderMove,
	}

	return cmd
}

func (c *cli) runfolderMove(cmd *cobra.Command, args []string) error {
	return c.vc.FolderMove(context.Background(), args[0], args[1])
}

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderMoveUse     = "move <source folder> <destination folder>"
	folderMoveShort   = "Move a folder from source to destination"
	folderMoveExample = "vaku folder move secret/foo secret/bar"
	folderMoveLong    = "Move a folder from source to destination"
)

func (c *cli) newFolderMoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderMoveUse,
		Short:   folderMoveShort,
		Long:    folderMoveLong,
		Example: folderMoveExample,

		Args: cobra.ExactArgs(2), //nolint:gomnd

		RunE: c.runfolderMove,
	}

	return cmd
}

func (c *cli) runfolderMove(cmd *cobra.Command, args []string) error {
	return c.vc.FolderMove(context.Background(), args[0], args[1])
}

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderReadUse     = "read <path>"
	folderReadShort   = "Recursively read all paths in a folder"
	folderReadExample = "vaku folder read secret/foo"
	folderReadLong    = "Recursively read all paths in a folder"
)

func (c *cli) newFolderReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderReadUse,
		Short:   folderReadShort,
		Long:    folderReadLong,
		Example: folderReadExample,

		Args: cobra.ExactArgs(1),

		RunE: c.runfolderRead,
	}

	return cmd
}

func (c *cli) runfolderRead(cmd *cobra.Command, args []string) error {
	read, err := c.vc.FolderRead(context.Background(), args[0])
	c.output(read)
	return err
}

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderReadArgs    = 1
	folderReadUse     = "read <folder>"
	folderReadShort   = "Recursively read all secrets in a folder"
	folderReadLong    = "Recursively read all secrets in a folder"
	folderReadExample = "vaku folder read secret/foo"
)

func (c *cli) newFolderReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderReadUse,
		Short:   folderReadShort,
		Long:    folderReadLong,
		Example: folderReadExample,

		Args: cobra.ExactArgs(folderReadArgs),

		RunE: c.runfolderRead,
	}

	return cmd
}

func (c *cli) runfolderRead(cmd *cobra.Command, args []string) error {
	read, err := c.vc.FolderRead(context.Background(), args[0])
	c.output(read)
	return err
}

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderSearchUse     = "search <path> <search>"
	folderSearchShort   = "Search for a secret in a folder"
	folderSearchExample = "vaku folder search secret/foo"
	folderSearchLong    = "Search for a secret in a folder"
)

func (c *cli) newFolderSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderSearchUse,
		Short:   folderSearchShort,
		Long:    folderSearchLong,
		Example: folderSearchExample,

		Args: cobra.ExactArgs(2), //nolint:gomnd

		RunE: c.runfolderSearch,
	}

	return cmd
}

func (c *cli) runfolderSearch(cmd *cobra.Command, args []string) error {
	search, err := c.vc.FolderSearch(context.Background(), args[0], args[1])
	c.output(search)
	return err
}

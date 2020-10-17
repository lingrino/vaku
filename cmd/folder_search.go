package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

const (
	folderSearchArgs    = 2
	folderSearchUse     = "search <folder> <search>"
	folderSearchShort   = "Recursively search all secrets in a folder for a search string"
	folderSearchLong    = "Recursively search all secrets in a folder for a search string"
	folderSearchExample = "vaku folder search secret/foo bar"
)

func (c *cli) newFolderSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderSearchUse,
		Short:   folderSearchShort,
		Long:    folderSearchLong,
		Example: folderSearchExample,

		Args: cobra.ExactArgs(folderSearchArgs),

		RunE: c.runfolderSearch,
	}

	return cmd
}

func (c *cli) runfolderSearch(cmd *cobra.Command, args []string) error {
	search, err := c.vc.FolderSearch(context.Background(), args[0], args[1])
	c.output(search)
	return err
}

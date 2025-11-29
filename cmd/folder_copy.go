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

	cmd.Flags().Bool(flagAllVersionsName, flagAllVersionsDefault, flagAllVersionsUse)

	return cmd
}

func (c *cli) runfolderCopy(cmd *cobra.Command, args []string) error {
	allVersions, err := cmd.Flags().GetBool(flagAllVersionsName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	src, dst := args[0], args[1]

	if allVersions {
		return c.vc.FolderCopyAllVersions(ctx, src, dst)
	}
	return c.vc.FolderCopy(ctx, src, dst)
}

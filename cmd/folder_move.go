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

	cmd.Flags().Bool(flagAllVersionsName, flagAllVersionsDefault, flagAllVersionsUse)
	cmd.Flags().Bool(flagDestroyName, flagDestroyDefault, flagDestroyUse)

	return cmd
}

func (c *cli) runfolderMove(cmd *cobra.Command, args []string) error {
	allVersions, err := cmd.Flags().GetBool(flagAllVersionsName)
	if err != nil {
		return err
	}
	destroy, err := cmd.Flags().GetBool(flagDestroyName)
	if err != nil {
		return err
	}

	ctx := context.Background()
	src, dst := args[0], args[1]

	// --all-versions already destroys all versions at source
	if allVersions {
		return c.vc.FolderMoveAllVersions(ctx, src, dst)
	}

	// --destroy: copy all secrets, then permanently delete all versions
	if destroy {
		if err := c.vc.FolderCopy(ctx, src, dst); err != nil {
			return err
		}
		return c.vc.FolderDeleteMeta(ctx, src)
	}

	// default: copy all secrets, soft delete current versions
	return c.vc.FolderMove(ctx, src, dst)
}

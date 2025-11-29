package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathMoveArgs    = 2
	pathMoveUse     = "move <source path> <destination path>"
	pathMoveShort   = "Move a secret from a source path to a destination path"
	pathMoveLong    = "Move a secret from a source path to a destination path"
	pathMoveExample = "vaku path move secret/foo secret/bar"

	flagDestroyName    = "destroy"
	flagDestroyUse     = "permanently destroy all versions at source after copy (KV v2 only)"
	flagDestroyDefault = false
)

func (c *cli) newPathMoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathMoveUse,
		Short:   pathMoveShort,
		Long:    pathMoveLong,
		Example: pathMoveExample,

		Args: cobra.ExactArgs(pathMoveArgs),

		RunE: c.runPathMove,
	}

	cmd.Flags().Bool(flagAllVersionsName, flagAllVersionsDefault, "move all versions of the secret (KV v2 only)")
	cmd.Flags().Bool(flagDestroyName, flagDestroyDefault, flagDestroyUse)

	return cmd
}

func (c *cli) runPathMove(cmd *cobra.Command, args []string) error {
	allVersions, err := cmd.Flags().GetBool(flagAllVersionsName)
	if err != nil {
		return err
	}
	destroy, err := cmd.Flags().GetBool(flagDestroyName)
	if err != nil {
		return err
	}

	src, dst := args[0], args[1]

	// --all-versions already destroys all versions at source
	if allVersions {
		return c.vc.PathMoveAllVersions(src, dst)
	}

	// --destroy: copy current version, then permanently delete all versions
	if destroy {
		if err := c.vc.PathCopy(src, dst); err != nil {
			return err
		}
		return c.vc.PathDeleteMeta(src)
	}

	// default: copy current version, soft delete current version
	return c.vc.PathMove(src, dst)
}

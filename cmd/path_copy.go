package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathCopyArgs    = 2
	pathCopyUse     = "copy <source path> <destination path>"
	pathCopyShort   = "Copy a secret from a source path to a destination path"
	pathCopyLong    = "Copy a secret from a source path to a destination path"
	pathCopyExample = "vaku path copy secret/foo secret/bar"

	flagAllVersionsName    = "all-versions"
	flagAllVersionsUse     = "copy all versions of the secret (KV v2 only)"
	flagAllVersionsDefault = false
)

func (c *cli) newPathCopyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathCopyUse,
		Short:   pathCopyShort,
		Long:    pathCopyLong,
		Example: pathCopyExample,

		Args: cobra.ExactArgs(pathCopyArgs),

		RunE: c.runPathCopy,
	}

	cmd.Flags().Bool(flagAllVersionsName, flagAllVersionsDefault, flagAllVersionsUse)

	return cmd
}

func (c *cli) runPathCopy(cmd *cobra.Command, args []string) error {
	allVersions, _ := cmd.Flags().GetBool(flagAllVersionsName)
	if allVersions {
		return c.vc.PathCopyAllVersions(args[0], args[1])
	}
	return c.vc.PathCopy(args[0], args[1])
}

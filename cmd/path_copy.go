package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathCopyUse     = "copy <source path> <destination path>"
	pathCopyShort   = "Copy a secret from a source path to a destination path"
	pathCopyExample = "vaku path copy secret/foo secret/bar"
	pathCopyLong    = "Search a secret for a string"
)

func (c *cli) newPathCopyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathCopyUse,
		Short:   pathCopyShort,
		Long:    pathCopyLong,
		Example: pathCopyExample,

		Args: cobra.ExactArgs(2), //nolint:gomnd

		RunE: c.runPathCopy,
	}

	return cmd
}

func (c *cli) runPathCopy(cmd *cobra.Command, args []string) error {
	return c.vc.PathCopy(args[0], args[1])
}

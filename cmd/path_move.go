package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathMoveUse     = "move <source path> <destination path>"
	pathMoveShort   = "Move a secret from a source path to a destination path"
	pathMoveExample = "vaku path move secret/foo secret/bar"
	pathMoveLong    = "Search a secret for a string"
)

func (c *cli) newPathMoveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathMoveUse,
		Short:   pathMoveShort,
		Long:    pathMoveLong,
		Example: pathMoveExample,

		Args: cobra.ExactArgs(2), //nolint:gomnd

		RunE: c.runPathMove,
	}

	return cmd
}

func (c *cli) runPathMove(cmd *cobra.Command, args []string) error {
	return c.vc.PathMove(args[0], args[1])
}

package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathDeleteMetaArgs    = 1
	pathDeleteMetaUse     = "delete-meta <path>"
	pathDeleteMetaShort   = "Delete all secret metadata and versions at a path"
	pathDeleteMetaLong    = "Delete all secret metadata and versions at a path"
	pathDeleteMetaExample = "vaku path delete-meta secret/foo"
)

func (c *cli) newPathDeleteMetaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathDeleteMetaUse,
		Short:   pathDeleteMetaShort,
		Long:    pathDeleteMetaLong,
		Example: pathDeleteMetaExample,

		Args: cobra.ExactArgs(pathDeleteMetaArgs),

		RunE: c.runPathDeleteMeta,
	}

	return cmd
}

func (c *cli) runPathDeleteMeta(cmd *cobra.Command, args []string) error {
	return c.vc.PathDeleteMeta(args[0])
}

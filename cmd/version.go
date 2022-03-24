package cmd

import (
	"github.com/spf13/cobra"

	vaku "github.com/lingrino/vaku/v2/api"
)

const (
	versionUse     = "version"
	versionShort   = "Print vaku version"
	versionExample = "vaku version"
)

func (c *cli) newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     versionUse,
		Short:   versionShort,
		Example: versionExample,

		Args: cobra.NoArgs,

		DisableFlagsInUseLine: true,

		RunE: c.runVersion,
	}

	return cmd
}

func (c *cli) runVersion(cmd *cobra.Command, args []string) error {
	output := map[string]any{
		"CLI": c.version,
		"API": vaku.Version(),
	}
	c.output(output)
	return nil
}

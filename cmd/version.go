package cmd

import (
	"github.com/lingrino/vaku/vaku"
	"github.com/spf13/cobra"
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

		Run: c.runVersion,
	}

	return cmd
}

func (c *cli) runVersion(cmd *cobra.Command, args []string) {
	c.output("CLI: " + c.version)
	c.output("API: " + vaku.Version())
}

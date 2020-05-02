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

		RunE: c.runVersion,
	}

	return cmd
}

func (c *cli) runVersion(cmd *cobra.Command, args []string) error {
	output := map[string]interface{}{
		"CLI": c.version,
		"API": vaku.Version(),
	}
	c.output(output)
	return nil
}

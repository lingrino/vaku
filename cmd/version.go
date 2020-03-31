package cmd

import (
	"fmt"

	"github.com/lingrino/vaku/vaku"
	"github.com/spf13/cobra"
)

const (
	versionUse     = "version"
	versionShort   = "Print vaku version"
	versionExample = "vaku version"
)

func newVersionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:     versionUse,
		Short:   versionShort,
		Example: versionExample,

		Args: cobra.NoArgs,

		DisableFlagsInUseLine: true,

		Run: func(cmd *cobra.Command, args []string) {
			runVersion(version)
		},
	}

	return cmd
}

func runVersion(version string) {
	fmt.Println("CLI:", version)
	fmt.Println("API:", vaku.Version())
}

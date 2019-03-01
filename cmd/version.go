package cmd

import (
	"github.com/lingrino/vaku/vaku"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns the current Vaku CLI and API versions",

	Args: cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		print(map[string]interface{}{
			"CLI": version,
			"API": vaku.Version(),
		})
	},
}

func init() {
	VakuCmd.AddCommand(versionCmd)
}

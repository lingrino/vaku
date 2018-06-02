package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pathListCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List a vault key/value path",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])
		input.TrimPathPrefix = trimPathPrefix

		output, err := vgc.PathList(input)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to list path %s", args[0]))
		} else {
			print(map[string]interface{}{
				args[0]: output,
			})
		}
	},
}

func init() {
	pathCmd.AddCommand(pathListCmd)
	pathListCmd.PersistentFlags().BoolVarP(&trimPathPrefix, "trim-path-prefix", "t", true, "Output paths with the input path trimmed (like Vault CLI)")
}

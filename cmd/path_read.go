package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pathReadCmd = &cobra.Command{
	Use:   "read [path]",
	Short: "Read a vault path",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])
		input.TrimPathPrefix = trimPathPrefix

		output, err := vgc.PathRead(input)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to read path %s", args[0]))
		} else {
			print(map[string]interface{}{
				args[0]: output,
			})
		}
	},
}

func init() {
	pathCmd.AddCommand(pathReadCmd)
}

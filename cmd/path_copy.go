package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pathCopyCmd = &cobra.Command{
	Use:   "copy [source folder] [target path]",
	Short: "Copy a vault path from one location to another",

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		inputSource := vaku.NewPathInput(args[0])
		inputTarget := vaku.NewPathInput(args[1])

		err := vgc.PathCopy(inputSource, inputTarget)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to copy path %s to %s", args[0], args[1]))
		} else {
			print(map[string]interface{}{
				args[0]: fmt.Sprintf("Successfully copied path %s to %s", args[0], args[1]),
			})
		}
	},
}

func init() {
	pathCmd.AddCommand(pathCopyCmd)
}

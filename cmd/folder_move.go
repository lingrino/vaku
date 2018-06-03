package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderMoveCmd = &cobra.Command{
	Use:   "move [source folder] [target path]",
	Short: "Move a vault folder from one location to another",

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		inputSource := vaku.NewPathInput(args[0])
		inputTarget := vaku.NewPathInput(args[1])

		err := vgc.FolderMove(inputSource, inputTarget)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to move folder %s to %s", args[0], args[1]))
		} else {
			print(map[string]interface{}{
				args[0]: fmt.Sprintf("Successfully moved folder %s to %s", args[0], args[1]),
			})
		}
	},
}

func init() {
	folderCmd.AddCommand(folderMoveCmd)
}

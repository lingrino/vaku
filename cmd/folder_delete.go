package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderDeleteCmd = &cobra.Command{
	Use:   "delete [path]",
	Short: "Recursively delete an entire vault folder",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])

		err := vgc.FolderDelete(input)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to delete folder %s", args[0]))
		} else {
			print(map[string]interface{}{
				args[0]: "Successfully deleted folder, if it existed",
			})
		}
	},
}

func init() {
	folderCmd.AddCommand(folderDeleteCmd)
}

package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderDestroyCmd = &cobra.Command{
	Use:   "destroy [path]",
	Short: "Recursively destroy an entire vault folder (V2 mounts only)",
	Long: `Takes in a path and destroys every key in that folder and all sub-folders. Note that this
function only works on V2 mounts and that it destroys ALL versions of ALL keys

Example:
  vaku folder destroy secret/foo`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])

		err := vgc.FolderDestroy(input)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to destroy folder %s", args[0]))
		} else {
			print(map[string]interface{}{
				args[0]: "Successfully destroyed folder, if it existed",
			})
		}
	},
}

func init() {
	folderCmd.AddCommand(folderDestroyCmd)
}

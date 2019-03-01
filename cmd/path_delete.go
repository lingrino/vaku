package cmd

import (
	"fmt"

	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pathDeleteCmd = &cobra.Command{
	Use:   "delete [path]",
	Short: "Delete a vault path",
	Long: `Deletes a secret at a specified path. Note that for V2 mounts this only deletes the current version.
Functionally very similar to the 'vault delete' command, but works on v1 and v2 mounts.

Example:
  vaku path delete secret/foo`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])

		err := vgc.PathDelete(input)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to delete path %s", args[0]))
		} else {
			print(map[string]interface{}{
				args[0]: "Successfully deleted path, if it existed",
			})
		}
	},
}

func init() {
	pathCmd.AddCommand(pathDeleteCmd)
}

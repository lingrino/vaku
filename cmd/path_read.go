package cmd

import (
	"fmt"

	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pathReadCmd = &cobra.Command{
	Use:   "read [path]",
	Short: "Read a vault path",
	Long: `Reads a secret at a path. Functionally similar to 'vault read' but works on v1 and v2 mounts.

Example:
  vaku path read secret/foo`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])

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

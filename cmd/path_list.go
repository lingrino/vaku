package cmd

import (
	"fmt"

	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pathListCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "List a vault path",
	Long: `Lists all keys at a vault path. Functionally similar to 'vault list path' but works on v1 and v2 mounts.

Example:
  vaku path list secret/foo`,

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])
		input.TrimPathPrefix = !noTrimPathPrefix

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
	pathListCmd.Flags().BoolVarP(&noTrimPathPrefix, "no-trim-path-prefix", "T", false, "Output full paths instead of paths with the input path trimmed")
}

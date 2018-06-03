package cmd

import (
	"fmt"
	"os"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderListCmd = &cobra.Command{
	Use:   "list [path]",
	Short: "Recursively list a vault folder",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])
		input.TrimPathPrefix = !noTrimPathPrefix

		output, err := vgc.FolderList(input)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to list folder %s", args[0]))
			os.Exit(1)
		} else {
			print(map[string]interface{}{
				args[0]: output,
			})
		}
	},
}

func init() {
	folderCmd.AddCommand(folderListCmd)
	folderListCmd.PersistentFlags().BoolVarP(&noTrimPathPrefix, "no-trim-path-prefix", "T", true, "Output full paths instead of paths with the input path trimmed")
}

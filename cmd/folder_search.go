package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderSearchCmd = &cobra.Command{
	Use:   "search [folder] [search-string]",
	Short: "Search a vault folder for a string, returning all paths where it is found",

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])
		input.TrimPathPrefix = !noTrimPathPrefix

		output, err := vgc.FolderSearch(input, args[1])
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to search folder %s", args[0]))
		} else {
			print(map[string]interface{}{
				args[0]: output,
			})
		}
	},
}

func init() {
	folderCmd.AddCommand(folderSearchCmd)
	folderSearchCmd.PersistentFlags().BoolVarP(&noTrimPathPrefix, "no-trim-path-prefix", "T", true, "Output full paths instead of paths with the input path trimmed")
}

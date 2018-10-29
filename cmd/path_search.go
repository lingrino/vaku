package cmd

import (
	"fmt"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var pathSearchCmd = &cobra.Command{
	Use:   "search [path] [search-string]",
	Short: "Search a vault path for a string, returning true if it is found",
	Long: `Searches a vault secret at a path for a specified string. Note that this is a simple text search that
flattens the secret into a string and matches exactly the input provided. Returns true if the string is found and false
otherwise.

Example:
  vaku path search secret/foo "bar"`,

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])

		output, err := vgc.PathSearch(input, args[1])
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to search path %s", args[0]))
		} else {
			print(map[string]interface{}{
				args[0]: output,
			})
		}
	},
}

func init() {
	pathCmd.AddCommand(pathSearchCmd)
}

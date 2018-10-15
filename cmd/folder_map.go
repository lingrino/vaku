package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderMapCmd = &cobra.Command{
	Use:   "map [path]",
	Short: "Return a text map of the folder, with subfolders indented by depth",

	Args: cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		input := vaku.NewPathInput(args[0])
		input.TrimPathPrefix = true

		list, err := vgc.FolderList(input)
		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to list folder %s", args[0]))
			os.Exit(1)
		} else {
			var output string
			var prevPS []string
			var written bool

			// Loop over each return path
			for _, path := range list {
				// Split the path and loop over each piece of the path
				ps := strings.Split(path, "/")
				for psi, word := range ps {
					// Don't write anything if we've already written the "parent" word
					// Once we write one part of a path, we should write all of it
					if len(prevPS) > psi && word == prevPS[psi] && !written {
						continue
					}

					// Unless this is the last word, add a "/" to the output
					if len(ps) != psi+1 {
						word = word + "/"
					}
					output = output + strings.Repeat(indentString, psi) + word + "\n"
					written = true
				}
				prevPS = ps
				written = false
			}

			fmt.Println(output)
		}
	},
}

func init() {
	folderCmd.AddCommand(folderMapCmd)
	folderMapCmd.Flags().StringVarP(&indentString, "indent-string", "I", "    ", "The string to use for indenting the map")
}

package cmd

import (
	"fmt"

	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var folderCopyCmd = &cobra.Command{
	Use:   "copy [source folder] [target path]",
	Short: "Copy a vault folder from one location to another",
	Long: `Takes in a source path and target path and copies every path in the source to the target.
Note that this will copy the input path if it is a secret and all paths under the input path that
result from calling 'vaku folder list' on that path. Also note that this will overwrite any existing
keys at the target paths.

Example:
  vaku folder copy secret/foo secret/bar`,

	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		inputSource := vaku.NewPathInput(args[0])
		inputTarget := vaku.NewPathInput(args[1])

		var err error
		if !useSourceTarget {
			err = vgc.FolderCopy(inputSource, inputTarget)
		} else {
			authCopyClient()
			err = copyClient.FolderCopy(inputSource, inputTarget)
		}

		if err != nil {
			fmt.Printf("%s", errors.Wrapf(err, "Failed to copy folder %s to %s", args[0], args[1]))
		} else {
			print(map[string]interface{}{
				args[0]: fmt.Sprintf("Successfully copied folder %s to %s", args[0], args[1]),
			})
		}
	},
}

func init() {
	folderCmd.AddCommand(folderCopyCmd)
	folderCopyCmd.Flags().BoolVarP(&useSourceTarget, "use-source-target-params", "", false, "Use source|target address/namespace/token parameters")
	folderCopyCmd.Flags().StringVarP(&sourceAddress, "source-address", "", "", "The Vault address for the source folder")
	folderCopyCmd.Flags().StringVarP(&sourceNamespace, "source-namespace", "", "", "The Vault namespace for the source folder")
	folderCopyCmd.Flags().StringVarP(&sourceToken, "source-token", "", "", "The Vault token to access the source folder")
	folderCopyCmd.Flags().StringVarP(&targetAddress, "target-address", "", "", "The Vault address for the target folder")
	folderCopyCmd.Flags().StringVarP(&targetNamespace, "target-namespace", "", "", "The Vault namespace for the target folder")
	folderCopyCmd.Flags().StringVarP(&targetAddress, "target-token", "", "", "The Vault token to access the target folder")
}

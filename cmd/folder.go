package cmd

import (
	"github.com/spf13/cobra"
)

const (
	folderUse     = "folder <cmd>"
	folderShort   = "Commands that act on Vault folders"
	folderExample = "vaku folder list secret/foo"
	folderLong    = `Commands that act on Vault folders

Commands under the folder subcommand act on Vault folders. Folders
are designated by paths that end in a '/' such as 'secret/foo/'. Vaku
can list, copy, move, search, etc.. on Vault folders.`
)

func newFolderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderUse,
		Short:   folderShort,
		Long:    folderLong,
		Example: folderExample,
	}

	return cmd
}

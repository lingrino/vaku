package cmd

import (
	"github.com/spf13/cobra"
)

const (
	folderUse     = "folder <cmd>"
	folderShort   = "folder operations"
	folderExample = "vaku folder list secret/foo"
	folderLong    = `long
description
hello`
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

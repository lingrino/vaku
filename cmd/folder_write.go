package cmd

import (
	"context"
	"encoding/json"

	"github.com/spf13/cobra"
)

const (
	folderWriteArgs    = 1
	folderWriteUse     = "write <folder>"
	folderWriteShort   = "write a folder of secrets. WARNING: command expects a very specific json input"
	folderWriteLong    = "write a folder of secrets. WARNING: command expects a very specific json input"
	folderWriteExample = "vaku folder write '{\"a/b/c\": {\"foo\": \"bar\"}}'"
)

func (c *cli) newFolderWriteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     folderWriteUse,
		Short:   folderWriteShort,
		Long:    folderWriteLong,
		Example: folderWriteExample,

		Args: cobra.ExactArgs(folderWriteArgs),

		RunE: c.runfolderWrite,
	}

	return cmd
}

func (c *cli) runfolderWrite(cmd *cobra.Command, args []string) error {
	var input map[string]map[string]any
	err := json.Unmarshal([]byte(args[0]), &input)
	if err != nil {
		return c.combineErr(errJSONUnmarshal, err)
	}

	return c.vc.FolderWrite(context.Background(), input)
}

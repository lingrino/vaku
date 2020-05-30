package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathWriteUse     = "write <path>"
	pathWriteShort   = "Write all paths at a path"
	pathWriteExample = "vaku path write secret/foo"
	pathWriteLong    = "Write all paths at a path"
)

func (c *cli) newPathWriteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathWriteUse,
		Short:   pathWriteShort,
		Long:    pathWriteLong,
		Example: pathWriteExample,

		Args: cobra.ExactArgs(1),

		RunE: c.runPathWrite,
	}

	return cmd
}

func (c *cli) runPathWrite(cmd *cobra.Command, args []string) error {
	err := c.vc.PathWrite("kv2/data/wat", map[string]interface{}{"foo": args[0]})
	return err
}

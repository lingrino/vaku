package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathReadUse     = "read <path>"
	pathReadShort   = "Read all paths at a path"
	pathReadExample = "vaku path read secret/foo"
	pathReadLong    = "Read all paths at a path"
)

func (c *cli) newPathReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathReadUse,
		Short:   pathReadShort,
		Long:    pathReadLong,
		Example: pathReadExample,

		Args: cobra.ExactArgs(1),

		RunE: c.runPathRead,
	}

	return cmd
}

func (c *cli) runPathRead(cmd *cobra.Command, args []string) error {
	read, err := c.vc.PathRead(args[0])
	c.output(read)
	return err
}

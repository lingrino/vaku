package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathReadArgs    = 1
	pathReadUse     = "read <path>"
	pathReadShort   = "Read a secret at a path"
	pathReadLong    = "Read a secret at a path"
	pathReadExample = "vaku path read secret/foo"
)

func (c *cli) newPathReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathReadUse,
		Short:   pathReadShort,
		Long:    pathReadLong,
		Example: pathReadExample,

		Args: cobra.ExactArgs(pathReadArgs),

		RunE: c.runPathRead,
	}

	return cmd
}

func (c *cli) runPathRead(cmd *cobra.Command, args []string) error {
	read, err := c.vc.PathRead(args[0])
	c.output(read)
	return err
}

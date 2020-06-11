package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

const (
	pathSearchUse     = "search <path> <search>"
	pathSearchShort   = "Search a secret for a string"
	pathSearchExample = "vaku path search secret/foo bar"
	pathSearchLong    = "Search a secret for a string"
)

func (c *cli) newPathSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathSearchUse,
		Short:   pathSearchShort,
		Long:    pathSearchLong,
		Example: pathSearchExample,

		Args: cobra.ExactArgs(2), //nolint:gomnd

		RunE: c.runPathSearch,
	}

	return cmd
}

func (c *cli) runPathSearch(cmd *cobra.Command, args []string) error {
	found, err := c.vc.PathSearch(args[0], args[1])
	c.output(strconv.FormatBool(found))
	return err
}

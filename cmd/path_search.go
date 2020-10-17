package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
)

const (
	pathSearchArgs    = 2
	pathSearchUse     = "search <path> <search>"
	pathSearchShort   = "Search a secret for a search string"
	pathSearchLong    = "Search a secret for a search string"
	pathSearchExample = "vaku path search secret/foo bar"
)

func (c *cli) newPathSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathSearchUse,
		Short:   pathSearchShort,
		Long:    pathSearchLong,
		Example: pathSearchExample,

		Args: cobra.ExactArgs(pathSearchArgs),

		RunE: c.runPathSearch,
	}

	return cmd
}

func (c *cli) runPathSearch(cmd *cobra.Command, args []string) error {
	search, err := c.vc.PathSearch(args[0], args[1])
	c.output(strconv.FormatBool(search))
	return err
}

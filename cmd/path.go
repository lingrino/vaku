package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathUse     = "path [cmd]"
	pathShort   = "path operations"
	pathExample = "vaku path list secret/foo"
)

func newPathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathUse,
		Short:   pathShort,
		Example: pathExample,
	}

	return cmd
}

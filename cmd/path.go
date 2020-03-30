package cmd

import (
	"github.com/spf13/cobra"
)

const (
	pathUse     = "path <cmd>"
	pathShort   = "path operations"
	pathExample = "vaku path list secret/foo"
	pathLong    = `long
description
hello`
)

func newPathCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     pathUse,
		Short:   pathShort,
		Long:    pathLong,
		Example: pathExample,
	}

	return cmd
}

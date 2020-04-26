package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	docsArgs    = 1
	docsUse     = "docs <path>"
	docsShort   = "Generates markdown docs at a path"
	docsExample = "vaku docs ."
)

var (
	errDocGenMarkdown = errors.New("failed to generate markdown docs")
)

func (c *cli) newDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Hidden: true,

		Use:     docsUse,
		Short:   docsShort,
		Example: docsExample,

		Args: cobra.ExactArgs(docsArgs),

		DisableFlagsInUseLine: true,

		RunE: c.runDocs,
	}

	return cmd
}

func (c *cli) runDocs(cmd *cobra.Command, args []string) error {
	err := doc.GenMarkdownTree(cmd.Root(), args[0])
	if err != nil {
		return errDocGenMarkdown
	}

	return nil
}

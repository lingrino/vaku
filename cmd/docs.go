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
	errDocNilRoot     = errors.New("failed to generate docs for nil root command")
	errDocGenMarkdown = errors.New("failed to generate markdown docs")
)

func newDocsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Hidden: true,

		Use:     docsUse,
		Short:   docsShort,
		Example: docsExample,

		Args: cobra.ExactArgs(docsArgs),

		DisableFlagsInUseLine: true,

		RunE: func(cmd *cobra.Command, args []string) error {
			err := runDocs(cmd.Root(), args[0])
			return err
		},
	}

	return cmd
}

func runDocs(rootCmd *cobra.Command, folder string) error {
	if rootCmd == nil {
		return errDocNilRoot
	}

	err := doc.GenMarkdownTree(rootCmd, folder)
	if err != nil {
		return errDocGenMarkdown
	}

	return nil
}

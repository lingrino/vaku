package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

const (
	docsArgs    = 1
	docsUse     = "docs [path]"
	docsShort   = "generates markdown docs in the provided folder"
	docsExample = "vaku docs docs/cli"
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

		Run: func(cmd *cobra.Command, args []string) {
			err := runDocs(cmd.Root(), args[0])
			outErr(err)
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

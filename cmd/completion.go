package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

const (
	completionArgs    = 1
	completionUse     = "completion bash|zsh|powershell"
	completionShort   = "Generates shell completions"
	completionExample = "vaku completion zsh"
	completionLong    = `To install completions for your shell

# Bash: In ~/.bashrc
source <(vaku completion bash)

# Zsh: In ~/.zshhrc
source <(vaku completion zsh)

# Powershell
Write the contents of 'vaku completion powershell' and source them in your profile`
)

var (
	errCmpNilRoot     = errors.New("failed to print completions for nil root command")
	errCmpUnsupported = errors.New("unsupported completion type")
	errCmpFailed      = errors.New("failed to print completions")
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     completionUse,
		Short:   completionShort,
		Long:    completionLong,
		Example: completionExample,

		Args:      cobra.ExactArgs(completionArgs),
		ValidArgs: []string{"zsh", "bash", "powershhell"},

		RunE: func(cmd *cobra.Command, args []string) error {
			err := runCompletion(cmd.Root(), args[0])
			return err
		},
	}

	return cmd
}

func runCompletion(rootCmd *cobra.Command, completion string) error {
	if rootCmd == nil {
		return errCmpNilRoot
	}

	var err error
	switch completion {
	case "bash":
		err = rootCmd.GenBashCompletion(rootCmd.OutOrStdout())
	case "zsh":
		err = rootCmd.GenZshCompletion(rootCmd.OutOrStdout())
	case "powershell":
		err = rootCmd.GenPowerShellCompletion(rootCmd.OutOrStdout())
	case "fail":
		err = errors.New("failure injection")
	default:
		return errCmpUnsupported
	}

	if err != nil {
		return errCmpFailed
	}
	return nil
}

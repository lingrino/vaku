package cmd

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

const (
	completionArgs    = 1
	completionUse     = "completion bash|zsh|powershell"
	completionShort   = "generates completion scripts for bash, zsh, or powershell"
	completionExample = "vaku completion zsh"
	completionLong    = `To install completions for your shell

# Bash: In ~/.bashrc
source <(vaku completion bash)

# Zsh: In ~/.zshhrc
source <(vaku completion zsh)

# Powershell
Write the contents of 'vaku completion powershell' and source them in your profile
`
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

		Run: func(cmd *cobra.Command, args []string) {
			err := runCompletion(cmd.Root(), args[0])
			outErr(err)
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
		err = rootCmd.GenBashCompletion(os.Stdout)
	case "zsh":
		err = rootCmd.GenZshCompletion(os.Stdout)
	case "powershell":
		err = rootCmd.GenPowerShellCompletion(os.Stdout)
	default:
		return errCmpUnsupported
	}

	if err != nil {
		return errCmpFailed
	}
	return nil
}

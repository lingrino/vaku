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

# Fish: In ~/.config/fish/config.fish
vaku completion fish | source -

# Powershell
Write the contents of 'vaku completion powershell' and source them in your profile

# Zsh: In ~/.zshhrc
source <(vaku completion zsh)`
)

var (
	errCmpUnsupported = errors.New("unsupported completion type")
	errCmpFailed      = errors.New("failed to print completions")
)

func (c *cli) newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     completionUse,
		Short:   completionShort,
		Long:    completionLong,
		Example: completionExample,

		Args:      cobra.ExactArgs(completionArgs),
		ValidArgs: []string{"bash", "fish", "powershhell", "zsh"},

		DisableFlagsInUseLine: true,

		RunE: c.runCompletion,
	}

	return cmd
}

func (c *cli) runCompletion(cmd *cobra.Command, args []string) error {
	rootCmd := cmd.Root()

	var err error
	switch args[0] {
	case "bash":
		err = rootCmd.GenBashCompletion(rootCmd.OutOrStdout())
	case "fish":
		err = rootCmd.GenFishCompletion(rootCmd.OutOrStdout(), true)
	case "powershell":
		err = rootCmd.GenPowerShellCompletion(rootCmd.OutOrStdout())
	case "zsh":
		err = rootCmd.GenZshCompletion(rootCmd.OutOrStdout())
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

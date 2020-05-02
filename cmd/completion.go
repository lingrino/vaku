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
	outW := rootCmd.OutOrStdout()

	var err error
	var cmpErr error

	switch args[0] {
	case "bash":
		cmpErr = rootCmd.GenBashCompletion(outW)
	case "fish":
		cmpErr = rootCmd.GenFishCompletion(outW, true)
	case "powershell":
		cmpErr = rootCmd.GenPowerShellCompletion(outW)
	case "zsh":
		cmpErr = rootCmd.GenZshCompletion(outW)
	case "fail":
		cmpErr = errors.New("fault injection")
	default:
		err = errCmpUnsupported
	}

	if cmpErr != nil {
		err = errCmpFailed
	}

	return err
}

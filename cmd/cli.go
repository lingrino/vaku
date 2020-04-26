package cmd

import (
	"fmt"
	"io"

	"github.com/lingrino/vaku/vaku"
	"github.com/spf13/cobra"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

// cli extends cobra.Command with our own config.
type cli struct {
	// clients
	vc  *vaku.Client
	cmd *cobra.Command

	// flags
	flagAbsPath bool
	flagFormat  string
	flagWorkers int

	// data
	version string
}

// newCLI returns a new CLI ready to run. Vaku client is not set because some commands (version) do
// not need it. Instead vc is initialized as a persistent function on the path/folder subcommands.
func newCLI() *cli {
	cli := &cli{}
	cli.cmd = cli.newVakuCmd()
	return cli
}

// setVersion sets the CLI version.
func (c *cli) setVersion(version string) {
	c.version = version
}

// Execute runs the CLI.
func Execute(version string, args []string, outW, errW io.Writer) int {
	cli := newCLI()
	cli.setVersion(version)

	cli.cmd.SetArgs(args)
	cli.cmd.SetOut(outW)
	cli.cmd.SetErr(errW)
	err := cli.cmd.Execute()
	if err != nil {
		fmt.Fprintf(cli.cmd.ErrOrStderr(), "Error: %s\n", err)
		return exitFailure
	}

	return exitSuccess
}

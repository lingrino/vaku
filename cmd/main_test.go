package cmd

import (
	"bytes"
	"testing"
)

// newTestCLI returns a CLI ready for running tests.
func newTestCLI(t *testing.T, args []string) (*cli, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	cli := newCLI()
	cli.flagIndent = ""

	var outW, errW bytes.Buffer
	cli.cmd.SetOut(&outW)
	cli.cmd.SetErr(&errW)

	cli.cmd.SetArgs(args)

	return cli, &outW, &errW
}

// newTestCLI returns a CLI with an initialized API ready for running tests.
func newTestCLIWithAPI(t *testing.T, args []string) (*cli, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	cli, outW, errW := newTestCLI(t, args)
	// cli.
	return cli, outW, errW
}

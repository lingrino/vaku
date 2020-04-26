package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
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

// assertError checks an error against an expected string (or nil) in that error.
func assertError(t *testing.T, err error, contains string) {
	t.Helper()

	if contains == "" {
		assert.NoError(t, err)
	} else {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), contains)
	}
}

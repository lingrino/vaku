package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// prepCmd sets the command args/output to test values and returns the writer that will be used
func prepCmd(t *testing.T, cmd *cobra.Command, args []string) *bytes.Buffer {
	t.Helper()

	cmd.SetArgs(args)

	var b bytes.Buffer
	cmd.SetOut(&b)

	return &b
}

// assertError checks an error against an expected string (or nil) in that error
func assertError(t *testing.T, err error, contains string) {
	t.Helper()

	if contains == "" {
		assert.NoError(t, err)
	} else {
		assert.Contains(t, err.Error(), contains)
	}
}

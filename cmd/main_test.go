package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// prepCmd sets args/output to test values and returns stdin/stderr writers
func prepCmd(t *testing.T, cmd *cobra.Command, args []string) (*bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	cmd.SetArgs(args)

	var out, err bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&err)

	return &out, &err
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

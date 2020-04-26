package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	t.Parallel()

	cli := newCLI()
	assert.Nil(t, cli.vc)
	assert.NotNil(t, cli.cmd)
	assert.Equal(t, "", cli.version)

	cli.setVersion("1.0.0")
	assert.Equal(t, "1.0.0", cli.version)
}

func TestExecute(t *testing.T) {
	t.Parallel()

	var outW, errW bytes.Buffer

	code := Execute("dev", os.Args[1:], &outW, &errW)
	assert.Equal(t, exitSuccess, code)

	code = Execute("dev", []string{"INVALID"}, &outW, &errW)
	assert.Equal(t, exitFailure, code)
}

// TestHasExample tests that every command has an example.
func TestHasExample(t *testing.T) {
	cli, _, _ := newTestCLI(t, nil)
	assert.True(t, allHasExample(cli.cmd))
}

// allHasExample recursively checks a command and it's children for example functions.
func allHasExample(cmds ...*cobra.Command) bool {
	res := true
	for _, cmd := range cmds {
		res = res && cmd.HasExample() && allHasExample(cmd.Commands()...)
	}
	return res
}

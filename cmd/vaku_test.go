package cmd

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestVaku(t *testing.T) {
	t.Parallel()

	vc := newVakuCmd("test")
	stdO, stdE := prepCmd(t, vc, nil)
	assert.Equal(t, "", stdE.String())

	err := vc.Execute()
	assert.NoError(t, err)
	assert.Contains(t, stdO.String(), vakuShort)
	assert.Contains(t, stdO.String(), vakuLong)
	assert.Contains(t, stdO.String(), vakuExample)
}

func TestExecute(t *testing.T) {
	t.Parallel()

	code := Execute("test")
	assert.Equal(t, exitSuccess, code)
	code = Execute("fail")
	assert.Equal(t, exitFail, code)
}

// TestHasExample tests that every command has an example
func TestHasExample(t *testing.T) {
	rootCmd := newVakuCmd("")
	assert.True(t, allHasExample(rootCmd))
}

// allHasExample recursively checks a command and it's children for example functions
func allHasExample(cmds ...*cobra.Command) bool {
	res := true
	for _, cmd := range cmds {
		res = res && cmd.HasExample() && allHasExample(cmd.Commands()...)
	}
	return res
}

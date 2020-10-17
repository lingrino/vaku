package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVaku(t *testing.T) {
	t.Parallel()

	cli, outW, errW := newTestCLI(t, nil)
	assert.Equal(t, "", errW.String())

	err := cli.cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, outW.String(), vakuShort)
	assert.Contains(t, outW.String(), vakuLong)
	assert.Contains(t, outW.String(), vakuExample)
}

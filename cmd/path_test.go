package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPath(t *testing.T) {
	t.Parallel()

	cli, outW, errW := newTestCLI(t, []string{"path"})
	assert.Equal(t, "", errW.String())

	err := cli.cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, outW.String(), pathShort)
	assert.Contains(t, outW.String(), pathLong)
}

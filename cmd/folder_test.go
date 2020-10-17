package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolder(t *testing.T) {
	t.Parallel()

	cli, outW, errW := newTestCLI(t, []string{"folder"})
	assert.Equal(t, "", errW.String())

	err := cli.cmd.Execute()

	assert.NoError(t, err)
	assert.Contains(t, outW.String(), folderShort)
	assert.Contains(t, outW.String(), folderLong)
}

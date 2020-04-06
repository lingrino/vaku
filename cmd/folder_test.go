package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolder(t *testing.T) {
	t.Parallel()

	vc := newFolderCmd()
	stdO, stdE := prepCmd(t, vc, nil)
	assert.Equal(t, "", stdE.String())

	err := vc.Execute()
	assert.NoError(t, err)
	assert.Contains(t, stdO.String(), folderShort)
	assert.Contains(t, stdO.String(), folderLong)
}

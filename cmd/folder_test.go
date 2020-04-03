package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolder(t *testing.T) {
	t.Parallel()

	vc := newFolderCmd()
	out, _ := prepCmd(t, vc, nil)

	err := vc.Execute()
	assert.NoError(t, err)
	assert.Contains(t, out.String(), folderShort)
	assert.Contains(t, out.String(), folderLong)
}

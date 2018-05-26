package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFolderDeleteData struct {
	input     *PathInput
	outputErr bool
}

func TestFolderDelete(t *testing.T) {
	c := clientInitForTests(t)

	tests := map[int]TestFolderDeleteData{
		1: {
			input:     NewPathInput("secretv1/test"),
			outputErr: false,
		},
		2: {
			input:     NewPathInput("secretv2/test"),
			outputErr: false,
		},
		3: {
			input:     NewPathInput("secretdoesnotexist/test"),
			outputErr: true,
		},
	}

	for _, d := range tests {
		e := c.FolderDelete(d.input)
		_, re := c.FolderRead(d.input)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Error(t, re)
			assert.NoError(t, e)
		}
	}

	// Reseed the vault server after tests end
	seed(t, c)
}

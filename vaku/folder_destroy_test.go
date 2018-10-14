package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestFolderDestroyData struct {
	input     *vaku.PathInput
	outputErr bool
}

func TestFolderDestroy(t *testing.T) {
	c := clientInitForTests(t)
	defer seed(t, c)

	tests := map[int]TestFolderDestroyData{
		1: {
			input:     vaku.NewPathInput("secretv1/test"),
			outputErr: true,
		},
		2: {
			input:     vaku.NewPathInput("secretv2/test"),
			outputErr: false,
		},
		3: {
			input:     vaku.NewPathInput("secretdoesnotexist/test"),
			outputErr: true,
		},
	}

	for _, d := range tests {
		e := c.FolderDestroy(d.input)
		_, re := c.FolderRead(d.input)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Error(t, re)
			assert.NoError(t, e)
		}
	}
}

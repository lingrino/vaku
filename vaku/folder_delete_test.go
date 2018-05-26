package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"

	"github.com/stretchr/testify/assert"
)

type TestFolderDeleteData struct {
	input     *vaku.PathInput
	outputErr bool
}

func TestFolderDelete(t *testing.T) {
	c := clientInitForTests(t)
	defer seed(t, c)

	tests := map[int]TestFolderDeleteData{
		1: {
			input:     vaku.NewPathInput("secretv1/test"),
			outputErr: false,
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
		e := c.FolderDelete(d.input)
		_, re := c.FolderRead(d.input)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Error(t, re)
			assert.NoError(t, e)
		}
	}
}

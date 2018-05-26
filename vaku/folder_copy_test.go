package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFolderCopyData struct {
	inputSource *PathInput
	inputTarget *PathInput
	outputErr   bool
}

func TestFolderCopy(t *testing.T) {
	c := NewClient()
	c.SimpleInit()

	tests := map[int]TestFolderCopyData{
		1: {
			inputSource: NewPathInput("secretv1/test"),
			inputTarget: NewPathInput("secretv1/foldercopy"),
			outputErr:   false,
		},
		2: {
			inputSource: NewPathInput("secretv2/test"),
			inputTarget: NewPathInput("secretv2/foldercopy"),
			outputErr:   false,
		},
		3: {
			inputSource: NewPathInput("secretv1/test"),
			inputTarget: NewPathInput("secretv2/foldercopy"),
			outputErr:   false,
		},
		4: {
			inputSource: NewPathInput("secretv2/test"),
			inputTarget: NewPathInput("secretv1/foldercopy"),
			outputErr:   false,
		},
		5: {
			inputSource: NewPathInput("secretdoesnotexist/test"),
			inputTarget: NewPathInput("secretv1/test"),
			outputErr:   true,
		},
		6: {
			inputSource: NewPathInput("secretv1/test"),
			inputTarget: NewPathInput("secretdoesnotexist/test"),
			outputErr:   true,
		},
	}

	for _, d := range tests {
		e := c.FolderCopy(d.inputSource, d.inputTarget)
		sr, _ := c.FolderRead(d.inputSource)
		tr, _ := c.FolderRead(d.inputTarget)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Equal(t, sr, tr)
			assert.NoError(t, e)
		}
	}
}

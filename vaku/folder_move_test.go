package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFolderMoveData struct {
	inputSource *PathInput
	inputTarget *PathInput
	outputErr   bool
}

func TestFolderMove(t *testing.T) {
	c := clientInitForTests(t)

	tests := map[int]TestFolderMoveData{
		1: {
			inputSource: NewPathInput("secretv1/test/inner"),
			inputTarget: NewPathInput("secretv1/foldermove"),
			outputErr:   false,
		},
		2: {
			inputSource: NewPathInput("secretv2/test/inner"),
			inputTarget: NewPathInput("secretv2/foldermove/inner"),
			outputErr:   false,
		},
		3: {
			inputSource: NewPathInput("secretv1/test"),
			inputTarget: NewPathInput("secretv2/foldermove"),
			outputErr:   false,
		},
		4: {
			inputSource: NewPathInput("secretv2/test"),
			inputTarget: NewPathInput("secretv1/foldermove"),
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
		c.FolderDelete(d.inputTarget)
		bsr, _ := c.FolderRead(d.inputSource)
		e := c.FolderMove(d.inputSource, d.inputTarget)
		_, sre := c.FolderRead(d.inputSource)
		tr, _ := c.FolderRead(d.inputTarget)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Equal(t, bsr, tr)
			assert.Error(t, sre)
			assert.NoError(t, e)
		}
		// Reseed the vault server after each test
		seed(t, c)
	}
}

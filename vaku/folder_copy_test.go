package vaku_test

import (
	"testing"

	"vaku/vaku"

	"github.com/stretchr/testify/assert"
)

type TestFolderCopyData struct {
	inputSource *vaku.PathInput
	inputTarget *vaku.PathInput
	outputErr   bool
}

func TestFolderCopy(t *testing.T) {
	c := clientInitForTests(t)

	tests := map[int]TestFolderCopyData{
		1: {
			inputSource: vaku.NewPathInput("secretv1/test"),
			inputTarget: vaku.NewPathInput("secretv1/foldercopy"),
			outputErr:   false,
		},
		2: {
			inputSource: vaku.NewPathInput("secretv2/test"),
			inputTarget: vaku.NewPathInput("secretv2/foldercopy"),
			outputErr:   false,
		},
		3: {
			inputSource: vaku.NewPathInput("secretv1/test"),
			inputTarget: vaku.NewPathInput("secretv2/foldercopy"),
			outputErr:   false,
		},
		4: {
			inputSource: vaku.NewPathInput("secretv2/test"),
			inputTarget: vaku.NewPathInput("secretv1/foldercopy"),
			outputErr:   false,
		},
		5: {
			inputSource: vaku.NewPathInput("secretdoesnotexist/test"),
			inputTarget: vaku.NewPathInput("secretv1/test"),
			outputErr:   true,
		},
		6: {
			inputSource: vaku.NewPathInput("secretv1/test"),
			inputTarget: vaku.NewPathInput("secretdoesnotexist/test"),
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

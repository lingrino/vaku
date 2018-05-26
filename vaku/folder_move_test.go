package vaku_test

import (
	"testing"

	"vaku/vaku"

	"github.com/stretchr/testify/assert"
)

type TestFolderMoveData struct {
	inputSource *vaku.PathInput
	inputTarget *vaku.PathInput
	outputErr   bool
}

func TestFolderMove(t *testing.T) {
	c := clientInitForTests(t)

	tests := map[int]TestFolderMoveData{
		1: {
			inputSource: vaku.NewPathInput("secretv1/test/inner"),
			inputTarget: vaku.NewPathInput("secretv1/foldermove"),
			outputErr:   false,
		},
		2: {
			inputSource: vaku.NewPathInput("secretv2/test/inner"),
			inputTarget: vaku.NewPathInput("secretv2/foldermove/inner"),
			outputErr:   false,
		},
		3: {
			inputSource: vaku.NewPathInput("secretv1/test"),
			inputTarget: vaku.NewPathInput("secretv2/foldermove"),
			outputErr:   false,
		},
		4: {
			inputSource: vaku.NewPathInput("secretv2/test"),
			inputTarget: vaku.NewPathInput("secretv1/foldermove"),
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
		seed(t, c) // reseed every time for this test
	}
}

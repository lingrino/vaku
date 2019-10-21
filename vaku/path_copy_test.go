package vaku_test

import (
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestPathCopyData struct {
	inputSource *vaku.PathInput
	inputTarget *vaku.PathInput
	outputErr   bool
}

func TestPathCopy(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestPathCopyData{
		1: {
			inputSource: vaku.NewPathInput("secretv1/test/foo"),
			inputTarget: vaku.NewPathInput("secretv1/pathcopy/foo"),
			outputErr:   false,
		},
		2: {
			inputSource: vaku.NewPathInput("secretv2/test/foo"),
			inputTarget: vaku.NewPathInput("secretv2/pathcopy/foo"),
			outputErr:   false,
		},
		3: {
			inputSource: vaku.NewPathInput("secretv1/test/fizz"),
			inputTarget: vaku.NewPathInput("secretv2/pathcopy/fizz"),
			outputErr:   false,
		},
		4: {
			inputSource: vaku.NewPathInput("secretv2/test/fizz"),
			inputTarget: vaku.NewPathInput("secretv1/pathcopy/fizz"),
			outputErr:   false,
		},
		5: {
			inputSource: vaku.NewPathInput("secretdoesnotexist/test/foo"),
			inputTarget: vaku.NewPathInput("secretv1/test/foo"),
			outputErr:   true,
		},
		6: {
			inputSource: vaku.NewPathInput("secretv1/test/foo"),
			inputTarget: vaku.NewPathInput("secretdoesnotexist/foo"),
			outputErr:   true,
		},
	}

	for _, d := range tests {
		e := c.PathCopy(d.inputSource, d.inputTarget)
		sr, _ := c.PathRead(d.inputSource)
		tr, _ := c.PathRead(d.inputTarget)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Equal(t, sr, tr)
			assert.NoError(t, e)
		}
	}
}


func TestCopyClientPathCopy(t *testing.T) {
	t.Parallel()
	c := copyClientInitForTests(t)

	tests := map[int]TestPathCopyData{
		1: {
			inputSource: vaku.NewPathInput("secretv1/test/foo"),
			inputTarget: vaku.NewPathInput("secretv1/pathcopy/foo"),
			outputErr:   false,
		},
		2: {
			inputSource: vaku.NewPathInput("secretv2/test/foo"),
			inputTarget: vaku.NewPathInput("secretv2/pathcopy/foo"),
			outputErr:   false,
		},
		3: {
			inputSource: vaku.NewPathInput("secretv1/test/fizz"),
			inputTarget: vaku.NewPathInput("secretv2/pathcopy/fizz"),
			outputErr:   false,
		},
		4: {
			inputSource: vaku.NewPathInput("secretv2/test/fizz"),
			inputTarget: vaku.NewPathInput("secretv1/pathcopy/fizz"),
			outputErr:   false,
		},
		5: {
			inputSource: vaku.NewPathInput("secretdoesnotexist/test/foo"),
			inputTarget: vaku.NewPathInput("secretv1/test/foo"),
			outputErr:   true,
		},
		6: {
			inputSource: vaku.NewPathInput("secretv1/test/foo"),
			inputTarget: vaku.NewPathInput("secretdoesnotexist/foo"),
			outputErr:   true,
		},
	}

	for _, d := range tests {
		e := c.PathCopy(d.inputSource, d.inputTarget)
		sr, _ := c.Source.PathRead(d.inputSource)
		tr, _ := c.Target.PathRead(d.inputTarget)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Equal(t, sr, tr)
			assert.NoError(t, e)
		}
	}
}

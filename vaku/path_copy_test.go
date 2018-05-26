package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPathCopyData struct {
	inputSource *PathInput
	inputTarget *PathInput
	outputErr   bool
}

func TestPathCopy(t *testing.T) {
	c := NewClient()
	c.SimpleInit()

	tests := map[int]TestPathCopyData{
		1: {
			inputSource: NewPathInput("secretv1/test/foo"),
			inputTarget: NewPathInput("secretv1/pathcopy/foo"),
			outputErr:   false,
		},
		2: {
			inputSource: NewPathInput("secretv2/test/foo"),
			inputTarget: NewPathInput("secretv2/pathcopy/foo"),
			outputErr:   false,
		},
		3: {
			inputSource: NewPathInput("secretv1/test/fizz"),
			inputTarget: NewPathInput("secretv2/pathcopy/fizz"),
			outputErr:   false,
		},
		4: {
			inputSource: NewPathInput("secretv2/test/fizz"),
			inputTarget: NewPathInput("secretv1/pathcopy/fizz"),
			outputErr:   false,
		},
		5: {
			inputSource: NewPathInput("secretdoesnotexist/test/foo"),
			inputTarget: NewPathInput("secretv1/test/foo"),
			outputErr:   true,
		},
		6: {
			inputSource: NewPathInput("secretv1/test/foo"),
			inputTarget: NewPathInput("secretdoesnotexist/foo"),
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

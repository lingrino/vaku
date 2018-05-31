package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestPathMoveData struct {
	inputSource *vaku.PathInput
	inputTarget *vaku.PathInput
	outputErr   bool
}

func TestPathMove(t *testing.T) {
	c := clientInitForTests(t)
	defer seed(t, c)

	tests := map[int]TestPathMoveData{
		1: {
			inputSource: vaku.NewPathInput("secretv1/test/foo"),
			inputTarget: vaku.NewPathInput("secretv1/pathmove/foo"),
			outputErr:   false,
		},
		2: {
			inputSource: vaku.NewPathInput("secretv2/test/foo"),
			inputTarget: vaku.NewPathInput("secretv2/pathmove/foo"),
			outputErr:   false,
		},
		3: {
			inputSource: vaku.NewPathInput("secretv1/test/fizz"),
			inputTarget: vaku.NewPathInput("secretv2/pathmove/fizz"),
			outputErr:   false,
		},
		4: {
			inputSource: vaku.NewPathInput("secretv2/test/fizz"),
			inputTarget: vaku.NewPathInput("secretv1/pathmove/fizz"),
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
		bsr, _ := c.PathRead(d.inputSource)
		e := c.PathMove(d.inputSource, d.inputTarget)
		_, sre := c.PathRead(d.inputSource)
		tr, _ := c.PathRead(d.inputTarget)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Equal(t, bsr, tr)
			assert.Error(t, sre)
			assert.NoError(t, e)
		}
	}

	// Reseed the vault server after tests end
	seed(t, c)
}

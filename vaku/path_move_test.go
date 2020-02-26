package vaku_test

import (
	"fmt"
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestPathMoveData struct {
	inputSource *vaku.PathInput
	inputTarget *vaku.PathInput
	outputErr   bool
}

func TestPathMove(t *testing.T) {
	var err error

	c := clientInitForTests(t)

	defer func() {
		err = seed(t, c)
		if err != nil {
			t.Error(fmt.Errorf("failed to reseed: %w", err))
		}
	}()

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
		sr, sre := c.PathRead(d.inputSource)
		tr, _ := c.PathRead(d.inputTarget)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			if sre == nil {
				assert.Equal(t, "SECRET_HAS_BEEN_DELETED", sr["VAKU_STATUS"])
			} else {
				assert.Error(t, sre)
			}
			assert.Equal(t, bsr, tr)
			assert.NoError(t, e)
		}
	}

	// Reseed the vault server after tests end
	err = seed(t, c)
	if err != nil {
		t.Error(fmt.Errorf("failed to reseed: %w", err))
	}
}

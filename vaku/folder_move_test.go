package vaku_test

import (
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type TestFolderMoveData struct {
	inputSource *vaku.PathInput
	inputTarget *vaku.PathInput
	outputErr   bool
}

func TestFolderMove(t *testing.T) {
	var err error

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
		err = c.FolderDelete(d.inputTarget)
		if err != nil {
			if d.outputErr {
				assert.Error(t, err)
			} else {
				t.Error(errors.Wrapf(err, "Failed to delete folder %s", d.inputTarget.Path))
			}
		}
		bsr, _ := c.FolderRead(d.inputSource)
		e := c.FolderMove(d.inputSource, d.inputTarget)
		sr, sre := c.FolderRead(d.inputSource)
		tr, _ := c.FolderRead(d.inputTarget)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			if sre == nil {
				for _, data := range sr {
					assert.Equal(t, "SECRET_HAS_BEEN_DELETED", data["VAKU_STATUS"])
				}
			} else {
				assert.Error(t, sre)
			}
			assert.Equal(t, bsr, tr)
			assert.NoError(t, e)
		}
		err = seed(t, c) // reseed every time for this test
		if err != nil {
			t.Error(errors.Wrap(err, "Failed to reseed"))
		}
	}
}

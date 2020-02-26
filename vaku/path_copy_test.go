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

type TestPathCopyDeletedData struct {
	inputSource   *vaku.PathInput
	inputTarget   *vaku.PathInput
	copyErr       bool
	oppath        string
}

func TestCopyClientPathDeleted(t *testing.T) {
	c := copyClientInitForTests(t)

	tests := map[int]TestPathCopyDeletedData{
		1: {
			inputSource:   vaku.NewPathInput("secretv1/copydeleted/test"),
			inputTarget:   vaku.NewPathInput("secretv1/copydeleted/test"),
			copyErr:       true,
			oppath:        "delete",
		},
		2: {
			inputSource:   vaku.NewPathInput("secretv2/copydeleted/test"),
			inputTarget:   vaku.NewPathInput("secretv2/copydeleted/test"),
			copyErr:       false,
			oppath:        "delete",
		},
		3: {
			inputSource:   vaku.NewPathInput("secretv2/copydestroyed/test"),
			inputTarget:   vaku.NewPathInput("secretv2/copydestroyed/test"),
			copyErr:       true,
			oppath:        "destroy",
		},
	}

	for _, d := range tests {
		secret := map[string]interface{}{
			"Eg5ljS7t": "6F1B5nBg",
			"quqr32S5": "81iY4HAN",
			"r6R0JUzX": "rs1mCRB5",
		}

		err := c.Source.PathWrite(d.inputSource, secret)
		assert.NoError(t, err)

		if d.oppath == "delete" {
			err = c.Source.PathDelete(d.inputSource)
			assert.NoError(t, err)
		} else if d.oppath == "destroy" {
			err = c.Source.PathDestroy(d.inputSource)
			assert.NoError(t, err)
		}

		copyErr := c.PathCopy(d.inputSource, d.inputTarget)
		tr, readTargetErr := c.Target.PathRead(d.inputTarget)
		if d.copyErr {
			assert.Error(t, copyErr)
			assert.Error(t, readTargetErr)
		} else {
			assert.NoError(t, copyErr)
			assert.Nil(t, tr)
			assert.Error(t, readTargetErr)
		}
	}
}

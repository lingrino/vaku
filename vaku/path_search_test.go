package vaku_test

import (
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestPathSearchData struct {
	inputPath   *vaku.PathInput
	inputSearch string
	output      bool
	outputErr   bool
}

func TestPathSearch(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestPathSearchData{
		1: {
			inputPath:   vaku.NewPathInput("secretv1/test/foo"),
			inputSearch: "ba",
			output:      true,
			outputErr:   false,
		},
		2: {
			inputPath:   vaku.NewPathInput("secretv2/test/foo"),
			inputSearch: "ba",
			output:      true,
			outputErr:   false,
		},
		3: {
			inputPath:   vaku.NewPathInput("secretv1/test/foo"),
			inputSearch: "value",
			output:      true,
			outputErr:   false,
		},
		4: {
			inputPath:   vaku.NewPathInput("secretv2/test/foo"),
			inputSearch: "value",
			output:      true,
			outputErr:   false,
		},
		5: {
			inputPath:   vaku.NewPathInput("secretv1/test/inner/again/inner/UCrt6sZT"),
			inputSearch: "s1mCR",
			output:      true,
			outputErr:   false,
		},
		6: {
			inputPath:   vaku.NewPathInput("secretv2/test/inner/again/inner/UCrt6sZT"),
			inputSearch: "s1mCR",
			output:      true,
			outputErr:   false,
		},
		7: {
			inputPath:   vaku.NewPathInput("secretv1/test/foo"),
			inputSearch: "eiojfdss",
			output:      false,
			outputErr:   false,
		},
		8: {
			inputPath:   vaku.NewPathInput("secretv2/test/foo"),
			inputSearch: "eiojfdss",
			output:      false,
			outputErr:   false,
		},
		9: {
			inputPath:   vaku.NewPathInput("secretv1/doesnotexist"),
			inputSearch: "foo",
			output:      false,
			outputErr:   true,
		},
		10: {
			inputPath:   vaku.NewPathInput("secretv2/doesnotexist"),
			inputSearch: "foo",
			output:      false,
			outputErr:   true,
		},
		11: {
			inputPath:   vaku.NewPathInput("secretdoesnotexist/test/foo"),
			inputSearch: "foo",
			output:      false,
			outputErr:   true,
		},
	}

	for _, d := range tests {
		o, e := c.PathSearch(d.inputPath, d.inputSearch)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestFolderSearchData struct {
	inputPath   *vaku.PathInput
	inputSearch string
	output      []string
	outputErr   bool
}

func TestFolderSearch(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestFolderSearchData{
		1: {
			inputPath:   vaku.NewPathInput("secretv1/test"),
			inputSearch: "ba",
			output: []string{
				"fizz",
				"foo",
				"value",
			},
			outputErr: false,
		},
		2: {
			inputPath:   vaku.NewPathInput("secretv2/test"),
			inputSearch: "ba",
			output: []string{
				"fizz",
				"foo",
				"value",
			},
			outputErr: false,
		},
		3: {
			inputPath:   vaku.NewPathInput("secretv1/test/inner"),
			inputSearch: "rs1mCRB5",
			output: []string{
				"again/inner/UCrt6sZT",
			},
			outputErr: false,
		},
		4: {
			inputPath:   vaku.NewPathInput("secretv2/test/inner"),
			inputSearch: "rs1mCRB5",
			output: []string{
				"again/inner/UCrt6sZT",
			},
			outputErr: false,
		},
		5: {
			inputPath: &vaku.PathInput{
				Path:           "secretv1/test",
				TrimPathPrefix: false,
			},
			inputSearch: "VbJ",
			output: []string{
				"secretv1/test/HToOeKKD",
			},
			outputErr: false,
		},
		6: {
			inputPath: &vaku.PathInput{
				Path:           "secretv2/test",
				TrimPathPrefix: false,
			},
			inputSearch: "VbJ",
			output: []string{
				"secretv2/test/HToOeKKD",
			},
			outputErr: false,
		},
		7: {
			inputPath:   vaku.NewPathInput("secretv1/doesnotexist"),
			inputSearch: "foo",
			output:      nil,
			outputErr:   true,
		},
		8: {
			inputPath:   vaku.NewPathInput("secretv2/doesnotexist"),
			inputSearch: "foo",
			output:      nil,
			outputErr:   true,
		},
		9: {
			inputPath:   vaku.NewPathInput("secretdoesnotexist/test"),
			inputSearch: "foo",
			output:      nil,
			outputErr:   true,
		},
	}

	for _, d := range tests {
		o, e := c.FolderSearch(d.inputPath, d.inputSearch)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

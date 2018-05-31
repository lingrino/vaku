package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestFolderListData struct {
	input     *vaku.PathInput
	output    []string
	outputErr bool
}

func TestFolderList(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestFolderListData{
		1: {
			input:     vaku.NewPathInput("secretv1/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/A2xlzTfE", "inner/WKNC3muM", "inner/again/inner/UCrt6sZT", "value"},
			outputErr: false,
		},
		2: {
			input:     vaku.NewPathInput("secretv2/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/A2xlzTfE", "inner/WKNC3muM", "inner/again/inner/UCrt6sZT", "value"},
			outputErr: false,
		},
		3: {
			input: &vaku.PathInput{
				Path:           "secretv1/test/inner",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv1/test/inner/A2xlzTfE", "secretv1/test/inner/WKNC3muM", "secretv1/test/inner/again/inner/UCrt6sZT"},
			outputErr: false,
		},
		4: {
			input: &vaku.PathInput{
				Path:           "secretv2/test/inner",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv2/test/inner/A2xlzTfE", "secretv2/test/inner/WKNC3muM", "secretv2/test/inner/again/inner/UCrt6sZT"},
			outputErr: false,
		},
		5: {
			input:     vaku.NewPathInput("secretv1/test/inner/again/inner"),
			output:    []string{"UCrt6sZT"},
			outputErr: false,
		},
		6: {
			input:     vaku.NewPathInput("secretv2/test/inner/again/inner"),
			output:    []string{"UCrt6sZT"},
			outputErr: false,
		},
		7: {
			input:     vaku.NewPathInput("secretv1/doesnotexist"),
			output:    nil,
			outputErr: true,
		},
		8: {
			input:     vaku.NewPathInput("secretv2/doesnotexist"),
			output:    nil,
			outputErr: true,
		},
	}

	for _, d := range tests {
		o, e := c.FolderList(d.input)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFolderListData struct {
	input     *PathInput
	output    []string
	outputErr bool
}

func TestFolderList(t *testing.T) {
	c := NewClient()
	c.simpleInit()

	tests := map[int]TestFolderListData{
		1: TestFolderListData{
			input:     NewPathInput("secretv1/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/A2xlzTfE", "inner/WKNC3muM", "inner/again/inner/UCrt6sZT", "value"},
			outputErr: false,
		},
		2: TestFolderListData{
			input:     NewPathInput("secretv2/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/A2xlzTfE", "inner/WKNC3muM", "inner/again/inner/UCrt6sZT", "value"},
			outputErr: false,
		},
		3: TestFolderListData{
			input: &PathInput{
				Path:           "secretv1/test/inner",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv1/test/inner/A2xlzTfE", "secretv1/test/inner/WKNC3muM", "secretv1/test/inner/again/inner/UCrt6sZT"},
			outputErr: false,
		},
		4: TestFolderListData{
			input: &PathInput{
				Path:           "secretv2/test/inner",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv2/test/inner/A2xlzTfE", "secretv2/test/inner/WKNC3muM", "secretv2/test/inner/again/inner/UCrt6sZT"},
			outputErr: false,
		},
		5: TestFolderListData{
			input:     NewPathInput("secretv1/test/inner/again/inner"),
			output:    []string{"UCrt6sZT"},
			outputErr: false,
		},
		6: TestFolderListData{
			input:     NewPathInput("secretv2/test/inner/again/inner"),
			output:    []string{"UCrt6sZT"},
			outputErr: false,
		},
		7: TestFolderListData{
			input:     NewPathInput("secretv1/doesnotexist"),
			output:    nil,
			outputErr: true,
		},
		8: TestFolderListData{
			input:     NewPathInput("secretv2/doesnotexist"),
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

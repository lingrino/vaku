package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPathListData struct {
	input     *PathInput
	output    []string
	outputErr bool
}

func TestPathList(t *testing.T) {
	c := NewClient()
	c.simpleInit()

	tests := map[int]TestPathListData{
		1: TestPathListData{
			input:     NewPathInput("secretv1/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
			outputErr: false,
		},
		2: TestPathListData{
			input:     NewPathInput("secretv2/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
			outputErr: false,
		},
		3: TestPathListData{
			input: &PathInput{
				Path:           "secretv1/test/inner/again/",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv1/test/inner/again/inner/"},
			outputErr: false,
		},
		4: TestPathListData{
			input: &PathInput{
				Path:           "secretv2/test/inner/again/",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv2/test/inner/again/inner/"},
			outputErr: false,
		},
		5: TestPathListData{
			input:     NewPathInput("secretv1/doesnotexist"),
			output:    nil,
			outputErr: true,
		},
		6: TestPathListData{
			input:     NewPathInput("secretv2/doesnotexist"),
			output:    nil,
			outputErr: true,
		},
	}

	for _, d := range tests {
		o, e := c.PathList(d.input)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

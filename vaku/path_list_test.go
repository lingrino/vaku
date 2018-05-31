package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestPathListData struct {
	input     *vaku.PathInput
	output    []string
	outputErr bool
}

func TestPathList(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestPathListData{
		1: {
			input:     vaku.NewPathInput("secretv1/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
			outputErr: false,
		},
		2: {
			input:     vaku.NewPathInput("secretv2/test"),
			output:    []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
			outputErr: false,
		},
		3: {
			input: &vaku.PathInput{
				Path:           "secretv1/test/inner/again/",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv1/test/inner/again/inner/"},
			outputErr: false,
		},
		4: {
			input: &vaku.PathInput{
				Path:           "secretv2/test/inner/again/",
				TrimPathPrefix: false,
			},
			output:    []string{"secretv2/test/inner/again/inner/"},
			outputErr: false,
		},
		5: {
			input:     vaku.NewPathInput("secretv1/doesnotexist"),
			output:    nil,
			outputErr: true,
		},
		6: {
			input:     vaku.NewPathInput("secretv2/doesnotexist"),
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

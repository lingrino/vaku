package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFolderWriteData struct {
	input     map[string]map[string]interface{}
	outputErr bool
}

func TestFolderWrite(t *testing.T) {
	c := NewClient()
	c.SimpleInit()

	tests := map[int]TestFolderWriteData{
		1: {
			input: map[string]map[string]interface{}{
				"secretv1/writetest/foo": map[string]interface{}{
					"value": "bar",
				},
				"secretv1/writetest/bar": map[string]interface{}{
					"value": "foo",
				},
			},
			outputErr: false,
		},
		2: {
			input: map[string]map[string]interface{}{
				"secretv2/writetest/foo": map[string]interface{}{
					"value": "bar",
				},
				"secretv2/writetest/bar": map[string]interface{}{
					"value": "foo",
				},
			},
			outputErr: false,
		},
		3: {
			input: map[string]map[string]interface{}{
				"secretv1/writetest/foo": map[string]interface{}{
					"value": "bar",
				},
				"secretv1/writetest/bar": map[string]interface{}{
					"value": "foo",
				},
				"secretv1/writetesttwo/foo/": map[string]interface{}{
					"value": "bar",
				},
			},
			outputErr: false,
		},
		4: {
			input: map[string]map[string]interface{}{
				"secretv2/writetest/foo": map[string]interface{}{
					"value": "bar",
				},
				"secretv2/writetest/bar": map[string]interface{}{
					"value": "foo",
				},
				"secretv2/writetesttwo/foo/": map[string]interface{}{
					"value": "bar",
				},
			},
			outputErr: false,
		},
		5: {
			input: map[string]map[string]interface{}{
				"secretdoesnotexist/writetest/foo": map[string]interface{}{
					"value": "bar",
				},
				"secretdoesnotexist/writetest/bar": map[string]interface{}{
					"value": "foo",
				},
			},
			outputErr: true,
		},
	}

	for _, d := range tests {
		e := c.FolderWrite(d.input)
		readBack := make(map[string]map[string]interface{})
		for k := range d.input {
			rb, _ := c.PathRead(NewPathInput(k))
			readBack[k] = rb
		}
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Equal(t, readBack, d.input)
			assert.NoError(t, e)
		}
	}
}

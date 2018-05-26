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
	c := clientInitForTests(t)

	tests := map[int]TestFolderWriteData{
		1: {
			input: map[string]map[string]interface{}{
				"secretv1/writetest/foo": {
					"value": "bar",
				},
				"secretv1/writetest/bar": {
					"value": "foo",
				},
			},
			outputErr: false,
		},
		2: {
			input: map[string]map[string]interface{}{
				"secretv2/writetest/foo": {
					"value": "bar",
				},
				"secretv2/writetest/bar": {
					"value": "foo",
				},
			},
			outputErr: false,
		},
		3: {
			input: map[string]map[string]interface{}{
				"secretv1/writetest/foo": {
					"value": "bar",
				},
				"secretv1/writetest/bar": {
					"value": "foo",
				},
				"secretv1/writetesttwo/foo/": {
					"value": "bar",
				},
			},
			outputErr: false,
		},
		4: {
			input: map[string]map[string]interface{}{
				"secretv2/writetest/foo": {
					"value": "bar",
				},
				"secretv2/writetest/bar": {
					"value": "foo",
				},
				"secretv2/writetesttwo/foo/": {
					"value": "bar",
				},
			},
			outputErr: false,
		},
		5: {
			input: map[string]map[string]interface{}{
				"secretdoesnotexist/writetest/foo": {
					"value": "bar",
				},
				"secretdoesnotexist/writetest/bar": {
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

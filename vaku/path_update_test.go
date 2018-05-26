package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPathUpdateData struct {
	inputPath    *PathInput
	inputData    map[string]interface{}
	expectedData map[string]interface{}
	outputErr    bool
}

func TestPathUpdate(t *testing.T) {
	c := clientInitForTests(t)

	tests := map[int]TestPathUpdateData{
		1: {
			inputPath: NewPathInput("secretv1/test/foo"),
			inputData: map[string]interface{}{
				"value": "buzz",
			},
			expectedData: map[string]interface{}{
				"value": "buzz",
			},
			outputErr: false,
		},
		2: {
			inputPath: NewPathInput("secretv2/test/foo"),
			inputData: map[string]interface{}{
				"value": "buzz",
			},
			expectedData: map[string]interface{}{
				"value": "buzz",
			},
			outputErr: false,
		},
		3: {
			inputPath: NewPathInput("secretv1/test/fizz"),
			inputData: map[string]interface{}{
				"foo": "buzz",
				"new": "boo",
			},
			expectedData: map[string]interface{}{
				"fizz": "buzz",
				"foo":  "buzz",
				"new":  "boo",
			},
			outputErr: false,
		},
		4: {
			inputPath: NewPathInput("secretv2/test/fizz"),
			inputData: map[string]interface{}{
				"foo": "buzz",
				"new": "boo",
			},
			expectedData: map[string]interface{}{
				"fizz": "buzz",
				"foo":  "buzz",
				"new":  "boo",
			},
			outputErr: false,
		},
		5: {
			inputPath: NewPathInput("secretdoesnotexist/test/fizz"),
			inputData: map[string]interface{}{
				"foo": "buzz",
				"new": "boo",
			},
			expectedData: map[string]interface{}{},
			outputErr:    true,
		},
	}

	for _, d := range tests {
		e := c.PathUpdate(d.inputPath, d.inputData)
		readBack, re := c.PathRead(d.inputPath)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Equal(t, d.expectedData, readBack)
			assert.NoError(t, e)
			assert.NoError(t, re)
		}
	}

	// Reseed the vault server after tests end
	seed(t, c)
}

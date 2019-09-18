package vaku_test

import (
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type TestPathUpdateData struct {
	inputPath    *vaku.PathInput
	inputData    map[string]interface{}
	expectedData map[string]interface{}
	outputErr    bool
}

func TestPathUpdate(t *testing.T) {
	var err error

	c := clientInitForTests(t)

	defer func() {
		err = seed(t, c)
		if err != nil {
			t.Error(errors.Wrap(err, "Failed to reseed"))
		}
	}()

	tests := map[int]TestPathUpdateData{
		1: {
			inputPath: vaku.NewPathInput("secretv1/test/foo"),
			inputData: map[string]interface{}{
				"value": "buzz",
			},
			expectedData: map[string]interface{}{
				"value": "buzz",
			},
			outputErr: false,
		},
		2: {
			inputPath: vaku.NewPathInput("secretv2/test/foo"),
			inputData: map[string]interface{}{
				"value": "buzz",
			},
			expectedData: map[string]interface{}{
				"value": "buzz",
			},
			outputErr: false,
		},
		3: {
			inputPath: vaku.NewPathInput("secretv1/test/fizz"),
			inputData: map[string]interface{}{
				"foo":      "buzz",
				"vaku.new": "boo",
			},
			expectedData: map[string]interface{}{
				"fizz":     "buzz",
				"foo":      "buzz",
				"vaku.new": "boo",
			},
			outputErr: false,
		},
		4: {
			inputPath: vaku.NewPathInput("secretv2/test/fizz"),
			inputData: map[string]interface{}{
				"foo":      "buzz",
				"vaku.new": "boo",
			},
			expectedData: map[string]interface{}{
				"fizz":     "buzz",
				"foo":      "buzz",
				"vaku.new": "boo",
			},
			outputErr: false,
		},
		5: {
			inputPath: vaku.NewPathInput("secretdoesnotexist/test/fizz"),
			inputData: map[string]interface{}{
				"foo":      "buzz",
				"vaku.new": "boo",
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
	err = seed(t, c)
	if err != nil {
		t.Error(errors.Wrap(err, "Failed to reseed"))
	}
}

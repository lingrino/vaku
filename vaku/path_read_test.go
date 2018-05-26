package vaku_test

import (
	"testing"
	"vaku/vaku"

	"github.com/stretchr/testify/assert"
)

type TestPathReadData struct {
	input     *vaku.PathInput
	output    map[string]interface{}
	outputErr bool
}

func TestPathRead(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestPathReadData{
		1: {
			input: vaku.NewPathInput("secretv1/test/foo"),
			output: map[string]interface{}{
				"value": "bar",
			},
			outputErr: false,
		},
		2: {
			input: vaku.NewPathInput("secretv2/test/foo"),
			output: map[string]interface{}{
				"value": "bar",
			},
			outputErr: false,
		},
		3: {
			input: vaku.NewPathInput("secretv1/test/inner/again/inner/UCrt6sZT"),
			output: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
			outputErr: false,
		},
		4: {
			input: vaku.NewPathInput("secretv2/test/inner/again/inner/UCrt6sZT"),
			output: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
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
		o, e := c.PathRead(d.input)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

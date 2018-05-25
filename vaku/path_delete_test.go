package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPathDeleteData struct {
	input     *PathInput
	outputErr bool
}

func TestPathDelete(t *testing.T) {
	c := NewClient()
	c.SimpleInit()

	tests := map[int]TestPathDeleteData{
		1: {
			input:     NewPathInput("secretv1/test/foo"),
			outputErr: false,
		},
		2: {
			input:     NewPathInput("secretv2/test/foo"),
			outputErr: false,
		},
		3: {
			input:     NewPathInput("secretv1/doesnotexist"),
			outputErr: false,
		},
		4: {
			input:     NewPathInput("secretv2/doesnotexist"),
			outputErr: false,
		},
		5: {
			input:     NewPathInput("secretdoesnotexist/test/foo"),
			outputErr: true,
		},
	}

	for _, d := range tests {
		e := c.PathDelete(d.input)
		_, re := c.PathRead(d.input)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.Error(t, re)
			assert.NoError(t, e)
		}
	}

	// Reseed the vault server after tests end
	c.seed()
}

package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestFolderReadData struct {
	input     *PathInput
	output    map[string]map[string]interface{}
	outputErr bool
}

func TestFolderRead(t *testing.T) {
	c := NewClient()
	c.simpleInit()

	tests := map[int]TestFolderReadData{
		1: TestFolderReadData{
			input: NewPathInput("secretv1/test"),
			output: map[string]map[string]interface{}{
				"foo": map[string]interface{}{
					"value": "bar",
				},
				"value": map[string]interface{}{
					"fizz": "buzz",
					"foo":  "bar",
				},
				"fizz": map[string]interface{}{
					"fizz": "buzz",
					"foo":  "bar",
				},
				"HToOeKKD": map[string]interface{}{
					"3zqxVbJY": "TvOjGxvC",
				},
			},
			outputErr: false,
		},
		2: TestFolderReadData{
			input: NewPathInput("secretv2/test"),
			output: map[string]map[string]interface{}{
				"foo": map[string]interface{}{
					"value": "bar",
				},
				"value": map[string]interface{}{
					"fizz": "buzz",
					"foo":  "bar",
				},
				"fizz": map[string]interface{}{
					"fizz": "buzz",
					"foo":  "bar",
				},
				"HToOeKKD": map[string]interface{}{
					"3zqxVbJY": "TvOjGxvC",
				},
			},
			outputErr: false,
		},
		3: TestFolderReadData{
			input: &PathInput{
				Path:           "secretv1/test/inner/again/inner/",
				TrimPathPrefix: false,
			},
			output: map[string]map[string]interface{}{
				"secretv1/test/inner/again/inner/UCrt6sZT": map[string]interface{}{
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		4: TestFolderReadData{
			input: &PathInput{
				Path:           "secretv2/test/inner/again/inner/",
				TrimPathPrefix: false,
			},
			output: map[string]map[string]interface{}{
				"secretv2/test/inner/again/inner/UCrt6sZT": map[string]interface{}{
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		5: TestFolderReadData{
			input: &PathInput{
				Path:           "secretv1/doesnotexist",
				TrimPathPrefix: false,
			},
			output:    nil,
			outputErr: true,
		},
		6: TestFolderReadData{
			input: &PathInput{
				Path:           "secretv2/doesnotexist",
				TrimPathPrefix: false,
			},
			output:    nil,
			outputErr: true,
		},
	}

	for _, d := range tests {
		o, e := c.FolderRead(d.input)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

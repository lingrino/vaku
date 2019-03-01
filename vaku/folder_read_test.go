package vaku_test

import (
	"testing"

	"github.com/lingrino/vaku/vaku"
	"github.com/stretchr/testify/assert"
)

type TestFolderReadData struct {
	input     *vaku.PathInput
	output    map[string]map[string]interface{}
	outputErr bool
}

func TestFolderReadOnce(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestFolderReadData{
		1: {
			input: vaku.NewPathInput("secretv1/test"),
			output: map[string]map[string]interface{}{
				"foo": {
					"value": "bar",
				},
				"value": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"fizz": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"HToOeKKD": {
					"3zqxVbJY": "TvOjGxvC",
				},
			},
			outputErr: false,
		},
		2: {
			input: vaku.NewPathInput("secretv2/test"),
			output: map[string]map[string]interface{}{
				"foo": {
					"value": "bar",
				},
				"value": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"fizz": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"HToOeKKD": {
					"3zqxVbJY": "TvOjGxvC",
				},
			},
			outputErr: false,
		},
		3: {
			input: &vaku.PathInput{
				Path:           "secretv1/test/inner/again/inner/",
				TrimPathPrefix: false,
			},
			output: map[string]map[string]interface{}{
				"secretv1/test/inner/again/inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		4: {
			input: &vaku.PathInput{
				Path:           "secretv2/test/inner/again/inner/",
				TrimPathPrefix: false,
			},
			output: map[string]map[string]interface{}{
				"secretv2/test/inner/again/inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		5: {
			input: &vaku.PathInput{
				Path:           "secretv1/doesnotexist",
				TrimPathPrefix: false,
			},
			output:    nil,
			outputErr: true,
		},
		6: {
			input: &vaku.PathInput{
				Path:           "secretv2/doesnotexist",
				TrimPathPrefix: false,
			},
			output:    nil,
			outputErr: true,
		},
	}

	for _, d := range tests {
		o, e := c.FolderReadOnce(d.input)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

func TestFolderRead(t *testing.T) {
	c := clientInitForTests(t)

	tests := map[int]TestFolderReadData{
		1: {
			input: vaku.NewPathInput("secretv1/test"),
			output: map[string]map[string]interface{}{
				"foo": {
					"value": "bar",
				},
				"value": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"fizz": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"HToOeKKD": {
					"3zqxVbJY": "TvOjGxvC",
				},
				"inner/WKNC3muM": {
					"IY1C148K": "JxBfEt91",
					"iwVzPqbY": "0NH9GlR1",
				},
				"inner/A2xlzTfE": {
					"Eg5ljS7t": "BHRMKjj1",
					"quqr32S5": "pcidzSMW",
				},
				"inner/again/inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		2: {
			input: vaku.NewPathInput("secretv2/test"),
			output: map[string]map[string]interface{}{
				"foo": {
					"value": "bar",
				},
				"value": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"fizz": {
					"fizz": "buzz",
					"foo":  "bar",
				},
				"HToOeKKD": {
					"3zqxVbJY": "TvOjGxvC",
				},
				"inner/WKNC3muM": {
					"IY1C148K": "JxBfEt91",
					"iwVzPqbY": "0NH9GlR1",
				},
				"inner/A2xlzTfE": {
					"Eg5ljS7t": "BHRMKjj1",
					"quqr32S5": "pcidzSMW",
				},
				"inner/again/inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		3: {
			input: &vaku.PathInput{
				Path:           "secretv1/test/inner/again/inner/",
				TrimPathPrefix: false,
			},
			output: map[string]map[string]interface{}{
				"secretv1/test/inner/again/inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		4: {
			input: &vaku.PathInput{
				Path:           "secretv2/test/inner/again/inner/",
				TrimPathPrefix: false,
			},
			output: map[string]map[string]interface{}{
				"secretv2/test/inner/again/inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			outputErr: false,
		},
		5: {
			input: &vaku.PathInput{
				Path:           "secretv1/doesnotexist",
				TrimPathPrefix: false,
			},
			output:    nil,
			outputErr: true,
		},
		6: {
			input: &vaku.PathInput{
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

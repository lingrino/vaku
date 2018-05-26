package vaku_test

import (
	"testing"
	"vaku/vaku"

	"github.com/stretchr/testify/assert"
)

func TestKeyIsFolder(t *testing.T) {
	inputToOutput := map[string]bool{
		"/":       true,
		"a/":      true,
		"a/b/":    true,
		"":        false,
		"a":       false,
		"a/b":     false,
		"123/456": false,
	}

	c := vaku.NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.KeyIsFolder(i))
	}
}

func TestKeyJoin(t *testing.T) {
	outputToInput := map[string][]string{
		"/":       {"/"},
		"a/":      {"a/"},
		"b":       {"b", ""},
		"a/b/c":   {"a/b", "c"},
		"d/e/f":   {"d/e/", "/f"},
		"g/h/i/":  {"/g/h/", "/i/"},
		"j/k/l/m": {"/j/", "/k/l", "m"},
	}

	c := vaku.NewClient()
	for o, i := range outputToInput {
		assert.Equal(t, o, c.KeyJoin(i...))
	}
}

func TestPathJoin(t *testing.T) {
	outputToInput := map[string][]string{
		"":        {"/"},
		"a":       {"a/"},
		"b":       {"b", ""},
		"a/b/c":   {"a/b", "c"},
		"d/e/f":   {"d/e/", "/f"},
		"g/h/i":   {"/g/h/", "/i/"},
		"j/k/l/m": {"/j/", "/k/l", "m"},
	}

	c := vaku.NewClient()
	for o, i := range outputToInput {
		assert.Equal(t, o, c.PathJoin(i...))
	}
}

func TestKeyClean(t *testing.T) {
	inputToOutput := map[string]string{
		"":      "",
		"/":     "/",
		"a":     "a",
		"b/":    "b/",
		"/c":    "c",
		"/d/":   "d/",
		"/e/f/": "e/f/",
	}

	c := vaku.NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.KeyClean(i))
	}
}

func TestPathClean(t *testing.T) {
	inputToOutput := map[string]string{
		"":      "",
		"a":     "a",
		"b/":    "b",
		"/c":    "c",
		"/d/":   "d",
		"/e/f/": "e/f",
	}

	c := vaku.NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.PathClean(i))
	}
}

func TestKeyBase(t *testing.T) {
	inputToOutput := map[string]string{
		"":      "",
		"/":     "",
		"a":     "a",
		"b/":    "b/",
		"c/d":   "d",
		"/e/f/": "f/",
	}

	c := vaku.NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.KeyBase(i))
	}
}

func TestPathBase(t *testing.T) {
	inputToOutput := map[string]string{
		"":      "",
		"/":     "",
		"a":     "a",
		"b/":    "b",
		"c/d":   "d",
		"/e/f/": "f",
	}

	c := vaku.NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.PathBase(i))
	}
}

type TestSliceKeyPrefixData struct {
	inputSlice  []string
	inputPrefix string
	output      []string
}

func TestSliceAddKeyPrefix(t *testing.T) {
	tests := map[int]TestSliceKeyPrefixData{
		1: {
			inputSlice:  []string{"a"},
			inputPrefix: "b",
			output:      []string{"b/a"},
		},
		2: {
			inputSlice:  []string{"/c/d/e/"},
			inputPrefix: "/f/",
			output:      []string{"f/c/d/e/"},
		},
		3: {
			inputSlice:  []string{"/g/"},
			inputPrefix: "h",
			output:      []string{"h/g/"},
		},
		4: {
			inputSlice:  []string{"i/j"},
			inputPrefix: "i",
			output:      []string{"i/i/j"},
		},
	}

	c := vaku.NewClient()
	for _, d := range tests {
		c.SliceAddKeyPrefix(d.inputSlice, d.inputPrefix)
		assert.Equal(t, d.output, d.inputSlice)
	}
}

func TestSliceTrimKeyPrefix(t *testing.T) {
	tests := map[int]TestSliceKeyPrefixData{
		1: {
			inputSlice:  []string{"a"},
			inputPrefix: "b",
			output:      []string{"a"},
		},
		2: {
			inputSlice:  []string{"/c/d/e/"},
			inputPrefix: "/c/",
			output:      []string{"d/e/"},
		},
		3: {
			inputSlice:  []string{"f/g"},
			inputPrefix: "f",
			output:      []string{"g"},
		},
		4: {
			inputSlice:  []string{"i/j"},
			inputPrefix: "k",
			output:      []string{"i/j"},
		},
	}

	c := vaku.NewClient()
	for _, d := range tests {
		c.SliceTrimKeyPrefix(d.inputSlice, d.inputPrefix)
		assert.Equal(t, d.output, d.inputSlice)
	}
}

package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathIsFolder(t *testing.T) {
	inputToOutput := map[string]bool{
		"/":       true,
		"a/":      true,
		"a/b/":    true,
		"":        false,
		"a":       false,
		"a/b":     false,
		"123/456": false,
	}

	c := NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.PathIsFolder(i))
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

	c := NewClient()
	for o, i := range outputToInput {
		assert.Equal(t, o, c.PathJoin(i...))
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

	c := NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.PathClean(i))
	}
}

func TestPathLength(t *testing.T) {
	inputToOutput := map[string]int{
		"":        0,
		"a":       1,
		"b/c":     2,
		"d/e/":    2,
		"/f/g/h/": 3,
	}

	c := NewClient()
	for i, o := range inputToOutput {
		assert.Equal(t, o, c.PathLength(i))
	}
}

func TestPathGetPrefix(t *testing.T) {
	outputToInput := map[string]map[string]interface{}{
		"a/b":   {"path": "a/b/c", "depth": 2},
		"d/e/f": {"path": "d/e/f/g", "depth": 3},
		"d":     {"path": "d/e/f/g/h", "depth": 1},
		"h/i":   {"path": "h/i", "depth": 0},
		"j/k":   {"path": "j/k", "depth": -1},
		"l/m":   {"path": "l/m", "depth": 10},
	}

	c := NewClient()
	for o, i := range outputToInput {
		assert.Equal(t, o, c.PathGetPrefix(i["path"].(string), i["depth"].(int)))
	}
}

func TestPathGetSuffix(t *testing.T) {
	outputToInput := map[string]map[string]interface{}{
		"b/c":   {"path": "a/b/c", "depth": 2},
		"e/f/g": {"path": "d/e/f/g", "depth": 3},
		"h":     {"path": "d/e/f/g/h", "depth": 1},
		"h/i":   {"path": "h/i", "depth": 0},
		"j/k":   {"path": "j/k", "depth": -1},
		"l/m":   {"path": "l/m", "depth": 10},
	}

	c := NewClient()
	for o, i := range outputToInput {
		assert.Equal(t, o, c.PathGetSuffix(i["path"].(string), i["depth"].(int)))
	}
}

func TestPathRemovePrefix(t *testing.T) {
	outputToInput := map[string]map[string]interface{}{
		"c":       {"path": "a/b/c", "depth": 2},
		"g":       {"path": "d/e/f/g", "depth": 3},
		"e/f/g/h": {"path": "d/e/f/g/h", "depth": 1},
		"h/i":     {"path": "h/i", "depth": 0},
		"j/k":     {"path": "j/k", "depth": -1},
		"l/m":     {"path": "l/m", "depth": 10},
	}

	c := NewClient()
	for o, i := range outputToInput {
		assert.Equal(t, o, c.PathRemovePrefix(i["path"].(string), i["depth"].(int)))
	}
}

func TestPathRemoveSuffix(t *testing.T) {
	outputToInput := map[string]map[string]interface{}{
		"a":       {"path": "a/b/c", "depth": 2},
		"d":       {"path": "d/e/f/g", "depth": 3},
		"d/e/f/g": {"path": "d/e/f/g/h", "depth": 1},
		"h/i":     {"path": "h/i", "depth": 0},
		"j/k":     {"path": "j/k", "depth": -1},
		"l/m":     {"path": "l/m", "depth": 10},
	}

	c := NewClient()
	for o, i := range outputToInput {
		assert.Equal(t, o, c.PathRemoveSuffix(i["path"].(string), i["depth"].(int)))
	}
}

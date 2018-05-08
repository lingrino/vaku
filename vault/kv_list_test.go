package vault

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// kvListTest holds the test data for KVList
// The 'Path' here should be specified without
// the leading mounts secretv1/2. Same with output
type kvListTest struct {
	input  *KVListInput
	output []string
}

func TestKVList(t *testing.T) {
	c := NewClient()
	c.simpleInit()
	mounts := []string{"secretv1", "secretv2"}

	nameToTest := map[string]*kvListTest{
		"1": &kvListTest{
			input: &KVListInput{
				Path:           "test",
				Recurse:        false,
				TrimPathPrefix: true,
			},
			output: []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
		},
		"2": &kvListTest{
			input: &KVListInput{
				Path:           "test",
				Recurse:        true,
				TrimPathPrefix: true,
			},
			output: []string{"HToOeKKD", "fizz", "foo", "inner/A2xlzTfE", "inner/WKNC3muM", "inner/again/inner/UCrt6sZT", "value"},
		},
		"3": &kvListTest{
			input: &KVListInput{
				Path:           "test/fizz",
				Recurse:        false,
				TrimPathPrefix: false,
			},
			output: nil,
		},
		"4": &kvListTest{
			input: &KVListInput{
				Path:           "test/inner",
				Recurse:        true,
				TrimPathPrefix: false,
			},
			output: []string{"test/inner/A2xlzTfE", "test/inner/WKNC3muM", "test/inner/again/inner/UCrt6sZT"},
		},
	}

	for _, mount := range mounts {
		for _, v := range nameToTest {
			var origPath string

			origPath = v.input.Path
			v.input.Path = c.PathJoin(mount, v.input.Path)

			l, _ := c.KVList(v.input)
			for i, p := range l {
				l[i] = strings.TrimPrefix(p, mount+"/")
			}

			assert.Equal(t, v.output, l)
			v.input.Path = origPath
		}
	}
}

package vault

import (
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
		// "3": &kvListTest{
		// 	input: &KVListInput{
		// 		Path:           "test",
		// 		Recurse:        true,
		// 		TrimPathPrefix: false,
		// 	},
		// 	output: []string{"test/HToOeKKD", "test/fizz", "test/foo", "test/inner/A2xlzTfE", "test/inner/WKNC3muM", "test/inner/again/inner/UCrt6sZT", "test/value"},
		// },
		// "4": &kvListTest{
		// 	input: &KVListInput{
		// 		Path:           "secretv1/test",
		// 		Recurse:        false,
		// 		TrimPathPrefix: true,
		// 	},
		// 	output: []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
		// },
		// "5": &kvListTest{
		// 	input: &KVListInput{
		// 		Path:           "secretv1/test",
		// 		Recurse:        false,
		// 		TrimPathPrefix: true,
		// 	},
		// 	output: []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
		// },
		// "6": &kvListTest{
		// 	input: &KVListInput{
		// 		Path:           "secretv1/test",
		// 		Recurse:        false,
		// 		TrimPathPrefix: true,
		// 	},
		// 	output: []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
		// },
		// "7": &kvListTest{
		// 	input: &KVListInput{
		// 		Path:           "secretv1/test",
		// 		Recurse:        false,
		// 		TrimPathPrefix: true,
		// 	},
		// 	output: []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
		// },
		// "8": &kvListTest{
		// 	input: &KVListInput{
		// 		Path:           "secretv1/test",
		// 		Recurse:        false,
		// 		TrimPathPrefix: true,
		// 	},
		// 	output: []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
		// },
	}

	for _, mount := range mounts {
		for _, v := range nameToTest {
			var shortPath string
			// var shortOutput []string

			// shortOutput = v.output

			shortPath = v.input.Path
			v.input.Path = c.PathJoin(mount, v.input.Path)

			if !v.input.TrimPathPrefix {
				for i, p := range v.output {
					v.output[i] = c.PathJoin(mount, p)
				}
			}

			l, _ := c.KVList(v.input)
			assert.Equal(t, v.output, l)
			v.input.Path = shortPath
		}
	}
}

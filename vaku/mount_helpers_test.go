package vaku_test

import (
	"testing"

	"github.com/Lingrino/vaku/vaku"

	"github.com/stretchr/testify/assert"
)

type TestMountInfoData struct {
	input     string
	output    *vaku.MountInfoOutput
	outputErr bool
}

func TestMountInfo(t *testing.T) {
	t.Parallel()
	c := clientInitForTests(t)

	tests := map[int]TestMountInfoData{
		1: {
			input: "secretv1/test",
			output: &vaku.MountInfoOutput{
				FullPath:      "secretv1/test",
				MountPath:     "secretv1",
				MountlessPath: "test",
				MountVersion:  "1",
			},
			outputErr: false,
		},
		2: {
			input: "secretv2/test",
			output: &vaku.MountInfoOutput{
				FullPath:      "secretv2/test",
				MountPath:     "secretv2",
				MountlessPath: "test",
				MountVersion:  "2",
			},
			outputErr: false,
		},
		3: {
			input: "secretv1/doesnotexist",
			output: &vaku.MountInfoOutput{
				FullPath:      "secretv1/doesnotexist",
				MountPath:     "secretv1",
				MountlessPath: "doesnotexist",
				MountVersion:  "1",
			},
			outputErr: false,
		},
		4: {
			input: "secretv2/doesnotexist",
			output: &vaku.MountInfoOutput{
				FullPath:      "secretv2/doesnotexist",
				MountPath:     "secretv2",
				MountlessPath: "doesnotexist",
				MountVersion:  "2",
			},
			outputErr: false,
		},
		5: {
			input:     "doesnotexist/test",
			output:    nil,
			outputErr: true,
		},
	}

	for _, d := range tests {
		o, e := c.MountInfo(d.input)
		assert.Equal(t, d.output, o)
		if d.outputErr {
			assert.Error(t, e)
		} else {
			assert.NoError(t, e)
		}
	}
}

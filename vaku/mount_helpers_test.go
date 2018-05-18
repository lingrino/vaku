package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestMountInfoData struct {
	input     string
	output    *MountInfoOutput
	outputErr bool
}

func TestMountInfo(t *testing.T) {
	c := NewClient()
	c.SimpleInit()

	tests := map[int]TestMountInfoData{
		1: TestMountInfoData{
			input: "secretv1/test",
			output: &MountInfoOutput{
				FullPath:      "secretv1/test",
				MountPath:     "secretv1",
				MountlessPath: "test",
				MountVersion:  "1",
			},
			outputErr: false,
		},
		2: TestMountInfoData{
			input: "secretv2/test",
			output: &MountInfoOutput{
				FullPath:      "secretv2/test",
				MountPath:     "secretv2",
				MountlessPath: "test",
				MountVersion:  "2",
			},
			outputErr: false,
		},
		3: TestMountInfoData{
			input: "secretv1/doesnotexist",
			output: &MountInfoOutput{
				FullPath:      "secretv1/doesnotexist",
				MountPath:     "secretv1",
				MountlessPath: "doesnotexist",
				MountVersion:  "1",
			},
			outputErr: false,
		},
		4: TestMountInfoData{
			input: "secretv2/doesnotexist",
			output: &MountInfoOutput{
				FullPath:      "secretv2/doesnotexist",
				MountPath:     "secretv2",
				MountlessPath: "doesnotexist",
				MountVersion:  "2",
			},
			outputErr: false,
		},
		5: TestMountInfoData{
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

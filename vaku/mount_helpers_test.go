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
				mountPath:     "secretv1",
				MountlessPath: "test",
				mountVersion:  "1",
			},
			outputErr: false,
		},
		2: TestMountInfoData{
			input: "secretv2/test",
			output: &MountInfoOutput{
				FullPath:      "secretv2/test",
				mountPath:     "secretv2",
				MountlessPath: "test",
				mountVersion:  "2",
			},
			outputErr: false,
		},
		3: TestMountInfoData{
			input: "secretv1/doesnotexist",
			output: &MountInfoOutput{
				FullPath:      "secretv1/doesnotexist",
				mountPath:     "secretv1",
				MountlessPath: "doesnotexist",
				mountVersion:  "1",
			},
			outputErr: false,
		},
		4: TestMountInfoData{
			input: "secretv2/doesnotexist",
			output: &MountInfoOutput{
				FullPath:      "secretv2/doesnotexist",
				mountPath:     "secretv2",
				MountlessPath: "doesnotexist",
				mountVersion:  "2",
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

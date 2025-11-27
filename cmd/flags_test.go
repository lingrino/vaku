package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateVakuFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give *cli
		want error
	}{
		{
			name: "valid",
			give: &cli{
				flagAbsPath:     true,
				flagNoAccessErr: true,
				flagFormat:      "text",
				flagIndent:      "----",
				flagSort:        true,
				flagWorkers:     100,
			},
			want: nil,
		},
		{
			name: "invalid format",
			give: &cli{
				flagAbsPath:     true,
				flagNoAccessErr: true,
				flagFormat:      "invalid",
				flagIndent:      "----",
				flagSort:        true,
				flagWorkers:     100,
			},
			want: errFlagInvalidFormat,
		},
		{
			name: "invalid workers",
			give: &cli{
				flagAbsPath:     true,
				flagNoAccessErr: true,
				flagFormat:      "text",
				flagIndent:      "----",
				flagSort:        true,
				flagWorkers:     0,
			},
			want: errFlagInvalidWorkers,
		},
		{
			name: "valid mount flags",
			give: &cli{
				flagFormat:       "text",
				flagWorkers:      10,
				flagMountPath:    "secret/",
				flagMountVersion: "2",
			},
			want: nil,
		},
		{
			name: "valid mount flags v1",
			give: &cli{
				flagFormat:       "text",
				flagWorkers:      10,
				flagMountPath:    "kv1/",
				flagMountVersion: "1",
			},
			want: nil,
		},
		{
			name: "invalid mount version",
			give: &cli{
				flagFormat:       "text",
				flagWorkers:      10,
				flagMountPath:    "secret/",
				flagMountVersion: "3",
			},
			want: errFlagInvalidMountVersion,
		},
		{
			name: "mount version without path",
			give: &cli{
				flagFormat:       "text",
				flagWorkers:      10,
				flagMountVersion: "1",
			},
			want: errFlagMountVersionNoPath,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.validateVakuFlags(nil, nil)
			if err != nil {
				assert.Equal(t, tt.want, err)
			}
		})
	}
}

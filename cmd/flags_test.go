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
		{
			name: "valid source mount flags",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagSrcMountPath:    "secret/",
				flagSrcMountVersion: "2",
			},
			want: nil,
		},
		{
			name: "valid source mount flags v1",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagSrcMountPath:    "kv1/",
				flagSrcMountVersion: "1",
			},
			want: nil,
		},
		{
			name: "invalid source mount version",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagSrcMountPath:    "secret/",
				flagSrcMountVersion: "3",
			},
			want: errFlagInvalidSrcMountVersion,
		},
		{
			name: "source mount version without path",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagSrcMountVersion: "1",
			},
			want: errFlagSrcMountVersionNoPath,
		},
		{
			name: "valid destination mount flags",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagDstMountPath:    "dest/",
				flagDstMountVersion: "2",
			},
			want: nil,
		},
		{
			name: "valid destination mount flags v1",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagDstMountPath:    "kv1/",
				flagDstMountVersion: "1",
			},
			want: nil,
		},
		{
			name: "invalid destination mount version",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagDstMountPath:    "dest/",
				flagDstMountVersion: "3",
			},
			want: errFlagInvalidDstMountVersion,
		},
		{
			name: "destination mount version without path",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagDstMountVersion: "1",
			},
			want: errFlagDstMountVersionNoPath,
		},
		{
			name: "valid both source and destination mount flags",
			give: &cli{
				flagFormat:          "text",
				flagWorkers:         10,
				flagSrcMountPath:    "src/",
				flagSrcMountVersion: "2",
				flagDstMountPath:    "dst/",
				flagDstMountVersion: "2",
			},
			want: nil,
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

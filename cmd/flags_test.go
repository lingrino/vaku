package cmd

import (
	"testing"
)

func TestValidateVakuFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give *cli
		want string
	}{
		{
			name: "valid",
			give: &cli{
				flagAbsPath: true,
				flagFormat:  "text",
				flagIndent:  "----",
				flagSort:    true,
				flagWorkers: 100,
			},
			want: "",
		},
		{
			name: "invalid format",
			give: &cli{
				flagAbsPath: true,
				flagFormat:  "invalid",
				flagIndent:  "----",
				flagSort:    true,
				flagWorkers: 100,
			},
			want: errFlagInvalidFormat.Error(),
		},
		{
			name: "invalid workers",
			give: &cli{
				flagAbsPath: true,
				flagFormat:  "text",
				flagIndent:  "----",
				flagSort:    true,
				flagWorkers: 0,
			},
			want: errFlagInvalidWorkers.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.give.validateVakuFlags(nil, nil)
			assertError(t, err, tt.want)
		})
	}
}

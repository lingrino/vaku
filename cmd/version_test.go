package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveVersion string
		giveArgs    []string
		wantOut     string
		wantErr     string
	}{
		{
			name:        "version test",
			giveVersion: "test",
			wantOut:     "API: 2.10.0\nCLI: test\n",
		},
		{
			name:        "version version",
			giveVersion: "version",
			wantOut:     "API: 2.10.0\nCLI: version\n",
		},
		{
			name:        "args",
			giveVersion: "version",
			giveArgs:    []string{"arg1", "arg2"},
			wantOut:     "",
			wantErr:     "unknown command",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{"version"}, tt.giveArgs...)
			cli, outW, errW := newTestCLI(t, args)
			cli.setVersion(tt.giveVersion)

			ec := cli.execute()
			assert.Equal(t, ec*len(errW.String()), len(errW.String()), "unexpected exit code")

			assert.Equal(t, tt.wantOut, outW.String())
		})
	}
}

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveArgs []string
		wantOut  string
		wantErr  string
	}{
		{
			name:     "success",
			giveArgs: []string{"/tmp"},
		},
		{
			name:     "failure",
			giveArgs: []string{"//\\#\\--%@&*/"},
			wantErr:  "ERROR: failed to generate markdown docs\n",
		},
		{
			name:     "extra args",
			giveArgs: []string{"///", "foo", "bar"},
			wantErr:  "ERROR: accepts 1 arg(s), received 3\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{"docs"}, tt.giveArgs...)
			cli, outW, errW := newTestCLI(t, args)

			ec := cli.execute()
			assert.Equal(t, ec*len(errW.String()), len(errW.String()), "unexpected exit code")

			assert.Equal(t, tt.wantOut, outW.String())
			assert.Equal(t, tt.wantErr, errW.String())
		})
	}
}

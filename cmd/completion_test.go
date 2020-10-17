package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompletion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveArgs []string
		wantErr  string
	}{
		{
			name:     "bash",
			giveArgs: []string{"bash"},
		},
		{
			name:     "fish",
			giveArgs: []string{"fish"},
		},
		{
			name:     "powershell",
			giveArgs: []string{"powershell"},
		},
		{
			name:     "zsh",
			giveArgs: []string{"zsh"},
		},
		{
			name:     "no args",
			giveArgs: []string{},
			wantErr:  "ERROR: accepts 1 arg(s), received 0\n",
		},
		{
			name:     "bad arg",
			giveArgs: []string{"badarg"},
			wantErr:  "ERROR: " + errCmpUnsupported.Error() + "\n",
		},
		{
			name:     "failure injection",
			giveArgs: []string{"fail"},
			wantErr:  "ERROR: " + errCmpFailed.Error() + "\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{"completion"}, tt.giveArgs...)
			cli, outW, errW := newTestCLI(t, args)

			ec := cli.execute()
			assert.Equal(t, ec*len(errW.String()), len(errW.String()), "unexpected exit code")

			assert.Equal(t, tt.wantErr, errW.String())
			if tt.wantErr == "" {
				// Expect a long string (the completion code)
				assert.Greater(t, len(outW.String()), 100)
			}
		})
	}
}

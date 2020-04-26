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
			wantErr:  "accepts 1 arg(s), received 0",
		},
		{
			name:     "bad arg",
			giveArgs: []string{"badarg"},
			wantErr:  errCmpUnsupported.Error(),
		},
		{
			name:     "failure injection",
			giveArgs: []string{"fail"},
			wantErr:  errCmpFailed.Error(),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{"completion"}, tt.giveArgs...)
			cli, outW, errW := newTestCLI(t, args)
			assert.Equal(t, "", errW.String())

			err := cli.cmd.Execute()

			assertError(t, err, tt.wantErr)
			if tt.wantErr == "" {
				// Expect a long string (the completion code)
				assert.Greater(t, len(outW.String()), 100)
			}
		})
	}
}

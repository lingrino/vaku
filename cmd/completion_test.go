package cmd

import (
	"errors"
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

			vc := newCompletionCmd()
			stdO, stdE := prepCmd(t, vc, tt.giveArgs)
			assert.Equal(t, "", stdE.String())

			err := vc.Execute()

			assertError(t, err, tt.wantErr)
			if tt.wantErr == "" {
				// Expect a long string
				assert.Greater(t, len(stdO.String()), 100)
			}
		})
	}
}

// TestRunCompletion explicitly with a nil command
func TestRunCompletion(t *testing.T) {
	err := runCompletion(nil, "")
	assert.True(t, errors.Is(err, errCmpNilRoot))
}

package cmd

import (
	"errors"
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
			wantErr:  "failed to generate markdown docs",
		},
		{
			name:     "extra args",
			giveArgs: []string{"///", "foo", "bar"},
			wantErr:  "accepts 1 arg(s), received 3",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vc := newDocsCmd()
			stdO, stdE := prepCmd(t, vc, tt.giveArgs)
			assert.Equal(t, "", stdE.String())

			err := vc.Execute()

			assertError(t, err, tt.wantErr)
			if tt.wantErr == "" {
				assert.Equal(t, tt.wantOut, stdO.String())
			}
		})
	}
}

// TestRunDocs explicitly with a nil command.
func TestRunDocs(t *testing.T) {
	err := runDocs(nil, "")
	assert.True(t, errors.Is(err, errDocNilRoot))
}

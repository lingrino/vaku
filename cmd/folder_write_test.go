package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveArgs []string
		wantRead string
		wantOut  string
		wantErr  string
	}{
		{
			name:     "empty",
			giveArgs: []string{"{}"},
			wantOut:  "",
			wantErr:  "",
		},
		{
			name:     "abcxyz",
			giveArgs: []string{"{\"a/b/c\": {\"foo\": \"bar\"}, \"x/y/z\": {\"bar\": \"foo\"}}"},
			wantOut:  "",
			wantErr:  "",
		},
		{
			name:     "error",
			giveArgs: []string{"error"},
			wantOut:  "",
			wantErr:  "ERROR: json unmarshal\ninvalid character 'e' looking for beginning of value\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{"folder", "write"}, tt.giveArgs...)
			cli, outW, errW := newTestCLIWithAPI(t, args)

			ec := cli.execute()
			assert.Equal(t, ec*len(errW.String()), len(errW.String()), "unexpected exit code")

			assert.Equal(t, tt.wantOut, outW.String())
			assert.Equal(t, tt.wantErr, errW.String())
		})
	}
}

package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		giveArgs []string
		wantOut  string
		wantErr  string
	}{
		{
			name:     "foo",
			giveArgs: []string{"foo", "bar"},
			wantOut:  "bim/bom\nfoo/bar\n",
			wantErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			args := append([]string{"folder", "search"}, tt.giveArgs...)
			cli, outW, errW := newTestCLIWithAPI(t, args)

			ec := cli.execute()
			assert.Equal(t, ec*len(errW.String()), len(errW.String()), "unexpected exit code")

			assert.Equal(t, tt.wantOut, outW.String())
			assert.Equal(t, tt.wantErr, errW.String())
		})
	}
}

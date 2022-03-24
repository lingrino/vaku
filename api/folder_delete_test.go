package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give         string
		wantReadBack map[string]map[string]any
		wantErr      []error
	}{
		{
			give:         "0/1",
			wantReadBack: nil,
			wantErr:      nil,
		},
		{
			give:         "0/4/13",
			wantReadBack: nil,
			wantErr:      nil,
		},
		{
			give:         "empty/path",
			wantReadBack: nil,
			wantErr:      nil,
		},
		{
			give: "0/4/13/24/25/error/list/inject",
			wantReadBack: map[string]map[string]any{
				"26/27": {
					"28": "29",
				},
			},
			wantErr: []error{ErrFolderDelete, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
		{
			give: "0/4/13/24/25/26/error/delete/inject",
			wantReadBack: map[string]map[string]any{
				"27": {
					"28": "29",
				},
			},
			wantErr: []error{ErrFolderDelete, ErrPathDelete, ErrVaultDelete},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				prefix := prefix
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.FolderDelete(context.Background(), PathJoin(prefix, tt.give))
					compareErrors(t, err, tt.wantErr)

					readBack, err := sharedVakuClean.FolderRead(context.Background(), PathJoin(prefix, tt.give))
					assert.NoError(t, err)
					assert.Equal(t, tt.wantReadBack, readBack)
				})
			}
		})
	}
}

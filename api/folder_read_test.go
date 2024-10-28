package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		want    map[string]map[string]any
		wantErr []error
	}{
		{
			give:    "0/1",
			want:    nil,
			wantErr: nil,
		},
		{
			give: "0/4/13/24/25",
			want: map[string]map[string]any{
				"26/27": {
					"28": "29",
				},
			},
			wantErr: nil,
		},
		{
			give: "0/4/13",
			want: map[string]map[string]any{
				"14": {
					"15": "16",
				},
				"17": {
					"18": "19",
					"20": "21",
					"22": "23",
				},
				"24/25/26/27": {
					"28": "29",
				},
			},
			wantErr: nil,
		},
		{
			give:    "error/list/inject",
			want:    nil,
			wantErr: []error{ErrFolderRead, ErrFolderReadChan, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
		{
			give:    "0/4/13/24/25/26/error/read/inject",
			want:    nil,
			wantErr: []error{ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					read, err := sharedVaku.FolderRead(context.Background(), PathJoin(prefix, tt.give))
					compareErrors(t, err, tt.wantErr)

					TrimPrefixMap(read, prefix)
					assert.Equal(t, tt.want, read)
				})
			}
		})
	}
}

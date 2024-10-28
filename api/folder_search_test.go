package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		give       string
		giveSearch string
		want       []string
		wantErr    []error
	}{
		{
			give:       "0",
			giveSearch: "notfound",
			want:       nil,
			wantErr:    nil,
		},
		{
			give:       "0/4/13/24",
			giveSearch: "7",
			want:       nil,
			wantErr:    nil,
		},
		{
			give:       "0/4/13",
			giveSearch: "3",
			want:       []string{"17"},
			wantErr:    nil,
		},
		{
			give:       "0/4",
			giveSearch: "2",
			want:       []string{"8", "13/17", "13/24/25/26/27"},
			wantErr:    nil,
		},
		{
			give:       "0/4/error/read/inject",
			giveSearch: "aaaaaaaaa",
			want:       nil,
			wantErr:    []error{ErrFolderSearch, ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
		},
		{
			give:       "0/4/funcdata/read/inject",
			giveSearch: "aaaaaaaaa",
			want:       nil,
			wantErr:    []error{ErrFolderSearch, ErrJSONMarshal},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					matches, err := sharedVaku.FolderSearch(context.Background(), PathJoin(prefix, tt.give), tt.giveSearch)
					compareErrors(t, err, tt.wantErr)

					TrimPrefixList(matches, prefix)
					assert.ElementsMatch(t, tt.want, matches)
				})
			}
		})
	}
}

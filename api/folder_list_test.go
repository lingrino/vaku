package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		want    []string
		wantErr []error
	}{
		{
			give:    "0/1",
			want:    nil,
			wantErr: nil,
		},
		{
			give: "0/4/",
			want: []string{
				"5",
				"8",
				"13/14",
				"13/17",
				"13/24/25/26/27",
			},
			wantErr: nil,
		},
		{
			give:    "error/list/inject",
			want:    []string{},
			wantErr: []error{ErrFolderList, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					list, err := sharedVaku.FolderList(context.Background(), PathJoin(prefix, tt.give))
					compareErrors(t, err, tt.wantErr)

					TrimPrefixList(list, prefix)
					assert.ElementsMatch(t, tt.want, list)
				})
			}
		})
	}
}

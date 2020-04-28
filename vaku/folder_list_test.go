package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveLogical logical
		giveOptions []Option
		want        []string
		wantErr     []error
	}{
		{
			name:    "test/foo",
			give:    "test/foo",
			want:    nil,
			wantErr: nil,
		},
		{
			name: "test/inner trailing slash",
			give: "test/inner/",
			want: []string{
				"WKNC3muM",
				"A2xlzTfE",
				"again/inner/UCrt6sZT",
			},
			wantErr: nil,
		},
		{
			name: "test/inner no slash and absolute path",
			give: "test/inner",
			want: []string{
				"test/inner/WKNC3muM",
				"test/inner/A2xlzTfE",
				"test/inner/again/inner/UCrt6sZT",
			},
			giveOptions: []Option{WithabsolutePath(true)},
			wantErr:     nil,
		},
		{
			name: "list error",
			give: "test/inner",
			want: []string{},
			giveLogical: &errLogical{
				err: errInject,
			},
			wantErr: []error{ErrFolderList, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, _ := testSetup(t, tt.giveLogical, nil, tt.giveOptions...)

			for _, ver := range kvMountVersions {
				ver := ver
				t.Run(ver, func(t *testing.T) {
					t.Parallel()

					path := addMountToPath(t, tt.give, ver)

					list, err := client.FolderList(context.Background(), path)
					compareErrors(t, err, tt.wantErr)

					TrimPrefixList(list, ver)
					assert.ElementsMatch(t, tt.want, list)
				})
			}
		})
	}
}

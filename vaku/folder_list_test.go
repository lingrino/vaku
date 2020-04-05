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
			giveOptions: []Option{WithAbsolutePath(true)},
			wantErr:     nil,
		},
		{
			name: "list error",
			give: "test/inner",
			want: []string{},
			giveLogical: &errLogical{
				err: errInject,
			},
			wantErr: []error{ErrFolderList, ErrPathList, ErrVaultList},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ln, client := testClient(t, tt.giveOptions...)
			defer ln.Close()
			updateLogical(t, client, tt.giveLogical, tt.giveLogical)

			for _, ver := range kvMountVersions {
				path := addMountToPath(t, tt.give, ver)

				list, err := client.FolderList(context.Background(), path)
				compareErrors(t, err, tt.wantErr)

				TrimListPrefix(list, ver)
				assert.ElementsMatch(t, tt.want, list)
			}
		})
	}
}

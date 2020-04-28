package vaku

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestFolderSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveSearch  string
		giveLogical logical
		giveOptions []Option
		want        []string
		wantErr     []error
	}{
		{
			name:       "path",
			give:       "test/foo",
			giveSearch: "ba",
			want:       nil,
			wantErr:    nil,
		},
		{
			name:       "folder",
			give:       "test/inner",
			giveSearch: "quqr32S5",
			want: []string{
				"A2xlzTfE",
				"again/inner/UCrt6sZT",
			},
			wantErr: nil,
		},
		{
			name:        "folder abs path",
			give:        "test/inner",
			giveSearch:  "quqr32S5",
			giveOptions: []Option{WithabsolutePath(true)},
			want: []string{
				"test/inner/A2xlzTfE",
				"test/inner/again/inner/UCrt6sZT",
			},
			wantErr: nil,
		},
		{
			name:        "all",
			give:        "test",
			giveSearch:  "b",
			giveOptions: []Option{WithabsolutePath(true)},
			want: []string{
				"test/foo",
				"test/value",
				"test/fizz",
				"test/HToOeKKD",
				"test/inner/WKNC3muM",
			},
			wantErr: nil,
		},
		{
			name: "read error",
			give: "test/inner/again",
			want: nil,
			giveLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr: []error{ErrFolderSearch, ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
		},
		{
			name:       "bad secret data",
			give:       "test/inner/again",
			giveSearch: "aaaaaaaaa",
			want:       nil,
			giveLogical: &errLogical{
				secret: &api.Secret{
					Data: map[string]interface{}{
						"Data": func() {},
					},
				},
				op: "Read",
			},
			wantErr: []error{ErrFolderSearch, ErrJSONMarshal},
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

					matches, err := client.FolderSearch(context.Background(), path, tt.giveSearch)
					compareErrors(t, err, tt.wantErr)

					TrimPrefixList(matches, ver)
					assert.ElementsMatch(t, tt.want, matches)
				})
			}
		})
	}
}

package vaku2

import (
	"errors"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestPathList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveLogical logical
		giveOptions []Option
		want        []string
		wantErr     error
		skipMount   bool
	}{
		{
			name:    "list test",
			give:    "test",
			want:    []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
			wantErr: nil,
		},
		{
			name:        "full path prefix",
			give:        "test/inner/again/",
			giveOptions: []Option{WithFullPath(true)},
			want:        []string{"test/inner/again/inner/"},
			wantErr:     nil,
		},
		{
			name:    "single secret",
			give:    "test/foo",
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "list bad path",
			give:    "doesnotexist",
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "no mount",
			give:    noMountPrefix,
			want:    nil,
			wantErr: nil,
		},
		{
			name: "list error",
			give: "test",
			giveLogical: &errLogical{
				err: errInject,
			},
			want:    nil,
			wantErr: ErrVaultList,
		},
		{
			name: "nil secret",
			give: "test",
			giveLogical: &errLogical{
				secret: nil,
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "nil data",
			give: "test",
			giveLogical: &errLogical{
				secret: &api.Secret{
					Data: nil,
				},
			},
			want:    nil,
			wantErr: nil,
		},
		{
			name: "no keys",
			give: "test",
			giveLogical: &errLogical{
				secret: &api.Secret{
					Data: map[string]interface{}{
						"notkeys": "notkeys",
					},
				},
			},
			want:    nil,
			wantErr: ErrDecodeSecret,
		},
		{
			name: "keys not []interface{}",
			give: "test",
			giveLogical: &errLogical{
				secret: &api.Secret{
					Data: map[string]interface{}{
						"keys": 1,
					},
				},
			},
			want:    nil,
			wantErr: ErrDecodeSecret,
		},
		{
			name: "keys not string",
			give: "test",
			giveLogical: &errLogical{
				secret: &api.Secret{
					Data: map[string]interface{}{
						"keys": []interface{}{
							1,
						},
					},
				},
			},
			want:    nil,
			wantErr: ErrDecodeSecret,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ln, client := testClient(t, tt.giveOptions...)
			defer ln.Close()

			updateLogical(t, client, tt.giveLogical)

			for _, ver := range kvMountVersions {
				path := addMountToPath(t, tt.give, ver)

				list, err := client.PathList(path)
				listD, errD := client.PathListDest(path)
				TrimListPrefix(list, ver)
				TrimListPrefix(listD, ver)

				assert.True(t, errors.Is(err, tt.wantErr))
				assert.True(t, errors.Is(errD, tt.wantErr))

				assert.Equal(t, tt.want, list)
				assert.Equal(t, tt.want, listD)
			}
		})
	}
}

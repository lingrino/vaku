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
		want        []string
		wantErr     error
	}{
		{
			name:    "list test",
			give:    "test",
			want:    []string{"HToOeKKD", "fizz", "foo", "inner/", "value"},
			wantErr: nil,
		},
		{
			name:    "list inner",
			give:    "test/inner/again/",
			want:    []string{"inner/"},
			wantErr: nil,
		},
		{
			name:    "list bad path",
			give:    "doesnotexist",
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

			ln, apiClient := testServer(t)
			defer ln.Close()

			client, err := NewClient(WithVaultClient(apiClient))
			assert.NoError(t, err)

			if tt.giveLogical != nil {
				client.sourceL = tt.giveLogical
			}

			for _, ver := range kvMountVersions {
				l, err := client.PathList(PathJoin(ver, tt.give))

				assert.True(t, errors.Is(err, tt.wantErr))
				assert.Equal(t, tt.want, l)
			}
		})
	}
}

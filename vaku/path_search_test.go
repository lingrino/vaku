package vaku

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestPathSearch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveSearch  string
		giveLogical logical
		giveOptions []Option
		wantSuccess bool
		wantErr     []error
	}{
		{
			name:        "test/foo key success",
			give:        "test/foo",
			giveSearch:  "bar",
			wantSuccess: true,
			wantErr:     nil,
		},
		{
			name:        "test/foo value success",
			give:        "test/foo",
			giveSearch:  "alue",
			wantSuccess: true,
			wantErr:     nil,
		},
		{
			name:        "test/foo fail",
			give:        "test/foo",
			giveSearch:  "valuebar",
			wantSuccess: false,
			wantErr:     nil,
		},
		{
			name:        "test/inner/again/inner/UCrt6sZT suuccess",
			give:        "test/inner/again/inner/UCrt6sZT",
			giveSearch:  "iY4HA",
			wantSuccess: true,
			wantErr:     nil,
		},
		{
			name:        "no path with string",
			give:        "pathdoesnotexist",
			giveSearch:  "searchstring",
			wantSuccess: false,
			wantErr:     nil,
		},
		{
			name:        "no path empty string",
			give:        "pathdoesnotexist",
			giveSearch:  "",
			wantSuccess: false,
			wantErr:     nil,
		},
		{
			name:        "bad mount",
			give:        noMountPrefix,
			giveSearch:  "searchstring",
			wantSuccess: false,
			wantErr:     nil,
		},
		{
			name:       "inject secret err",
			give:       "test/foo",
			giveSearch: "bar",
			giveLogical: &errLogical{
				err: errInject,
			},
			wantSuccess: false,
			wantErr:     []error{ErrPathSearch, ErrPathRead, ErrVaultRead},
		},
		{
			name:       "inject bad secret data",
			give:       "test/foo",
			giveSearch: "bar",
			giveLogical: &errLogical{
				secret: &api.Secret{
					Data: map[string]interface{}{
						"Data": func() {},
					},
				},
			},
			wantSuccess: false,
			wantErr:     []error{ErrPathSearch, ErrJSONMarshal},
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

					success, err := client.PathSearch(path, tt.giveSearch)

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.wantSuccess, success)
				})
			}
		})
	}
}

package vaku2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveLogical logical
		giveOptions []Option
		wantErr     []error
	}{
		{
			name:    "delete path",
			give:    "test/foo",
			wantErr: nil,
		},
		{
			name:    "nonexistent path",
			give:    "doesnotexist",
			wantErr: nil,
		},
		{
			name:    "no mount",
			give:    noMountPrefix,
			wantErr: []error{ErrPathDelete, ErrVaultDelete},
		},
		{
			name: "error",
			give: "delete/foo",
			giveLogical: &errLogical{
				err: errInject,
			},
			wantErr: []error{ErrPathDelete, ErrVaultDelete},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ln, client := testClient(t, tt.giveOptions...)
			defer ln.Close()
			readbackClient := cloneCLient(t, client)
			updateLogical(t, client, tt.giveLogical, tt.giveLogical)

			funcs := []func(string) error{
				client.PathDelete,
				client.PathDeleteDst,
			}

			for _, ver := range kvMountVersions {
				for _, f := range funcs {
					path := addMountToPath(t, tt.give, ver)

					err := f(path)
					compareErrors(t, err, tt.wantErr)

					readBack, err := readbackClient.PathRead(path)
					assert.NoError(t, err)
					assert.Nil(t, readBack)
				}
			}
		})
	}
}

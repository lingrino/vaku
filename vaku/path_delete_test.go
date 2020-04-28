package vaku

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

			client, rbClient := testSetup(t, tt.giveLogical, nil, tt.giveOptions...)

			for _, ver := range kvMountVersions {
				ver := ver
				t.Run(ver, func(t *testing.T) {
					t.Parallel()

					path := addMountToPath(t, tt.give, ver)

					err := client.PathDelete(path)
					compareErrors(t, err, tt.wantErr)

					readBack, err := rbClient.PathRead(path)
					assert.NoError(t, err)
					assert.Nil(t, readBack)
				})
			}
		})
	}
}

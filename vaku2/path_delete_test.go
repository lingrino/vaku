package vaku2

import (
	"errors"
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
		wantErr     error
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
			wantErr: ErrVaultDelete,
		},
		{
			name: "error",
			give: "delete/foo",
			giveLogical: &errLogical{
				err: errInject,
			},
			wantErr: ErrVaultDelete,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ln, client := testClient(t, tt.giveOptions...)
			defer ln.Close()
			readbackClient := cloneCLient(t, client)

			updateLogical(t, client, tt.giveLogical)

			for _, ver := range kvMountVersions {
				path := addMountToPath(t, tt.give, ver)

				err := client.PathDelete(path)
				errD := client.PathDeleteDest(path)
				assert.True(t, errors.Is(err, tt.wantErr), err)
				assert.True(t, errors.Is(errD, tt.wantErr), err)

				readBack, err := readbackClient.PathRead(path)
				readBackD, errD := readbackClient.PathReadDest(path)
				assert.NoError(t, err)
				assert.NoError(t, errD)
				assert.Nil(t, readBack)
				assert.Nil(t, readBackD)
			}
		})
	}
}

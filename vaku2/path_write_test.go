package vaku2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveData    map[string]interface{}
		giveLogical logical
		giveOptions []Option
		wantErr     error
	}{
		{
			name: "new path",
			give: "write/bar",
			giveData: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
			wantErr: nil,
		},
		{
			name: "overwrite",
			give: "test/foo",
			giveData: map[string]interface{}{
				"foo": "bar",
			},
			wantErr: nil,
		},
		{
			name:     "nil data",
			give:     "write/foo",
			giveData: nil,
			wantErr:  ErrVaultWrite,
		},
		{
			name:     "no mount",
			give:     noMountPrefix,
			giveData: nil,
			wantErr:  ErrVaultWrite,
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

				err := client.PathWrite(path, tt.giveData)
				errD := client.PathWriteDest(path, tt.giveData)
				assert.True(t, errors.Is(err, tt.wantErr), err)
				assert.True(t, errors.Is(errD, tt.wantErr), err)

				readBack, err := readbackClient.PathRead(path)
				readBackD, errD := readbackClient.PathReadDest(path)
				assert.NoError(t, err)
				assert.NoError(t, errD)
				assert.Equal(t, tt.giveData, readBack)
				assert.Equal(t, tt.giveData, readBackD)
			}
		})
	}
}

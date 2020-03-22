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

			ln, apiClient := testServer(t)
			defer ln.Close()

			client, err := NewClient(append(tt.giveOptions, WithVaultClient(apiClient))...)
			assert.NoError(t, err)

			backupL := client.sourceL
			if tt.giveLogical != nil {
				client.sourceL = tt.giveLogical
			}

			for _, ver := range kvMountVersions {
				path := tt.give
				if tt.give != noMountPrefix {
					path = PathJoin(ver, tt.give)
				}

				err := client.PathWrite(path, tt.giveData)
				assert.True(t, errors.Is(err, tt.wantErr))

				if tt.give == noMountPrefix {
					client.sourceL = backupL
					readBack, err := client.PathRead(PathJoin(ver, tt.give))
					if tt.giveLogical != nil {
						client.sourceL = tt.giveLogical
					}
					assert.NoError(t, err)

					assert.Equal(t, tt.giveData, readBack)
				}
			}
		})
	}
}

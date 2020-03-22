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

				err := client.PathDelete(path)
				assert.True(t, errors.Is(err, tt.wantErr))

				if tt.give == noMountPrefix {
					client.sourceL = backupL
					readBack, err := client.PathRead(PathJoin(ver, tt.give))
					if tt.giveLogical != nil {
						client.sourceL = tt.giveLogical
					}

					assert.NoError(t, err)
					assert.Nil(t, readBack)
				}
			}
		})
	}
}

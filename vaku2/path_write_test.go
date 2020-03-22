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
			name: "write/foo",
			give: "write/foo",
			giveData: map[string]interface{}{
				"foo": "bar",
			},
			wantErr: nil,
		},
		{
			name: "write/bar",
			give: "write/bar",
			giveData: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
			wantErr: nil,
		},
		{
			name:     "bad mount",
			give:     "mountdoesnotexist",
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

			if tt.giveLogical != nil {
				client.sourceL = tt.giveLogical
			}

			for _, ver := range kvMountVersions {
				err := client.PathWrite(PathJoin(ver, tt.give), tt.giveData)
				assert.True(t, errors.Is(err, tt.wantErr))

				readBack, err := client.PathRead(PathJoin(ver, tt.give))
				assert.NoError(t, err)

				assert.Equal(t, tt.giveData, readBack)
			}
		})
	}
}

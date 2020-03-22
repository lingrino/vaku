package vaku2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveLogical logical
		giveOptions []Option
		want        map[string]interface{}
		wantErr     error
	}{
		{
			name: "test/foo",
			give: "test/foo",
			want: map[string]interface{}{
				"value": "bar",
			},
			wantErr: nil,
		},
		{
			name: "test/inner/again/inner/UCrt6sZT",
			give: "test/inner/again/inner/UCrt6sZT",
			want: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
			wantErr: nil,
		},
		{
			name:    "doesnotexist",
			give:    "doesnotexist",
			want:    nil,
			wantErr: nil,
		},
		{
			name: "error",
			give: "test/foo",
			giveLogical: &errLogical{
				err: errInject,
			},
			want:    nil,
			wantErr: ErrVaultRead,
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
				read, err := client.PathRead(PathJoin(ver, tt.give))

				assert.True(t, errors.Is(err, tt.wantErr))
				assert.Equal(t, tt.want, read)
			}
		})
	}
}

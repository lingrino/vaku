package vaku

import (
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
		wantErr     []error
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
			name:    "no secret",
			give:    "doesnotexist",
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "no mount",
			give:    noMountPrefix,
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
			wantErr: []error{ErrPathRead, ErrVaultRead},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testClient(t, tt.giveOptions...)
			updateLogical(t, client, tt.giveLogical, nil)

			for _, ver := range kvMountVersions {
				path := addMountToPath(t, tt.give, ver)

				read, err := client.PathRead(path)

				compareErrors(t, err, tt.wantErr)
				assert.Equal(t, tt.want, read)
			}
		})
	}
}

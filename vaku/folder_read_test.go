package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveLogical logical
		giveOptions []Option
		want        map[string]map[string]interface{}
		wantErr     []error
	}{
		{
			name:    "test/foo",
			give:    "test/foo",
			want:    nil,
			wantErr: nil,
		},
		{
			name: "test/inner/again",
			give: "test/inner/again",
			want: map[string]map[string]interface{}{
				"inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			wantErr: nil,
		},
		{
			name: "test/inner/again asbolute path",
			give: "test/inner/again",
			want: map[string]map[string]interface{}{
				"test/inner/again/inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			giveOptions: []Option{WithabsolutePath(true)},
			wantErr:     nil,
		},
		{
			name: "test/inner/again list error",
			give: "test/inner/again",
			want: nil,
			giveLogical: &errLogical{
				err: errInject,
				op:  "List",
			},
			wantErr: []error{ErrFolderRead, ErrFolderReadChan, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
		{
			name: "test/inner/again read error",
			give: "test/inner/again",
			want: nil,
			giveLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr: []error{ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
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

					read, err := client.FolderRead(context.Background(), path)
					compareErrors(t, err, tt.wantErr)

					TrimPrefixMap(read, ver)
					assert.Equal(t, tt.want, read)
				})
			}
		})
	}
}

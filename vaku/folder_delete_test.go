package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		give         string
		giveLogical  logical
		giveOptions  []Option
		wantReadBack map[string]map[string]interface{}
		wantErr      []error
	}{
		{
			name:         "data path",
			give:         "test/foo",
			wantReadBack: nil,
			wantErr:      nil,
		},
		{
			name:         "nested data path",
			give:         "test/inner/again",
			wantReadBack: nil,
			wantErr:      nil,
		},
		{
			name:         "empty path",
			give:         "empty/path",
			wantReadBack: nil,
			wantErr:      nil,
		},
		{
			name: "list error",
			give: "test/inner/again",
			giveLogical: &errLogical{
				err: errInject,
				op:  "List",
			},
			wantReadBack: map[string]map[string]interface{}{
				"inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			wantErr: []error{ErrFolderDelete, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
		{
			name: "delete error",
			give: "test/inner/again",
			giveLogical: &errLogical{
				err: errInject,
				op:  "Delete",
			},
			wantReadBack: map[string]map[string]interface{}{
				"inner/UCrt6sZT": {
					"Eg5ljS7t": "6F1B5nBg",
					"quqr32S5": "81iY4HAN",
					"r6R0JUzX": "rs1mCRB5",
				},
			},
			wantErr: []error{ErrFolderDelete, ErrPathDelete, ErrVaultDelete},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := testClient(t, tt.giveOptions...)
			readbackClient := cloneCLient(t, client)
			updateLogical(t, client, tt.giveLogical, nil)

			for _, ver := range kvMountVersions {
				path := addMountToPath(t, tt.give, ver)

				err := client.FolderDelete(context.Background(), path)
				compareErrors(t, err, tt.wantErr)

				readBack, err := readbackClient.FolderRead(context.Background(), path)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantReadBack, readBack)
			}
		})
	}
}

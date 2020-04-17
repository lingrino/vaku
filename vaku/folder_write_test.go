package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		give         map[string]map[string]interface{}
		giveLogical  logical
		giveOptions  []Option
		wantReadBack map[string]map[string]interface{}
		wantErr      []error
	}{
		{
			name:         "nil",
			give:         nil,
			wantReadBack: nil,
			wantErr:      nil,
		},
		{
			name: "empty data",
			give: map[string]map[string]interface{}{
				"test/foo": nil,
			},
			wantReadBack: nil,
			wantErr:      []error{ErrFolderWrite, ErrFolderWriteChan, ErrPathWrite, ErrNilData},
		},
		{
			name: "overwrite",
			give: map[string]map[string]interface{}{
				"test/foo": {
					"bibim": "bap",
				},
			},
			wantReadBack: map[string]map[string]interface{}{
				"test/foo": {
					"bibim": "bap",
				},
			},
			wantErr: nil,
		},
		{
			name: "two new paths",
			give: map[string]map[string]interface{}{
				"new/boo": {
					"wat": "tot",
				},
				"new/boo/too": {
					"watoo": "totoo",
					"watee": "totee",
				},
			},
			wantReadBack: map[string]map[string]interface{}{
				"new/boo": {
					"wat": "tot",
				},
				"new/boo/too": {
					"watoo": "totoo",
					"watee": "totee",
				},
			},
			wantErr: nil,
		},
		{
			name: "two different paths",
			give: map[string]map[string]interface{}{
				"new/one/boo": {
					"wat": "tot",
				},
				"two/two/bootoo": {
					"watoo": "totoo",
					"watee": "totee",
				},
			},
			wantReadBack: map[string]map[string]interface{}{
				"new/one/boo": {
					"wat": "tot",
				},
				"two/two/bootoo": {
					"watoo": "totoo",
					"watee": "totee",
				},
			},
			wantErr: nil,
		},
		{
			name: "path write fail",
			give: map[string]map[string]interface{}{
				"failonwrite": {
					"foo": "bar",
				},
			},
			giveLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantReadBack: nil,
			wantErr:      []error{ErrFolderWrite, ErrFolderWriteChan, ErrPathWrite, ErrVaultWrite},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ln, client := testClient(t, tt.giveOptions...)
			defer ln.Close()
			readbackClient := cloneCLient(t, client)
			updateLogical(t, client, tt.giveLogical, tt.giveLogical)

			for _, ver := range kvMountVersions {
				writeMap := make(map[string]map[string]interface{}, len(tt.give))
				for path, data := range tt.give {
					writeMap[addMountToPath(t, path, ver)] = data
				}

				err := client.FolderWrite(context.Background(), writeMap)
				compareErrors(t, err, tt.wantErr)

				for path, data := range tt.wantReadBack {
					readBack, err := readbackClient.PathRead(addMountToPath(t, path, ver))
					assert.NoError(t, err)
					assert.Equal(t, data, readBack)
				}
			}
		})
	}
}

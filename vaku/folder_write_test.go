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
			wantErr:      []error{ErrFolderWrite, ErrPathWrite, ErrNilData},
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
				"failonwrite/error/write/inject": {
					"foo": "bar",
				},
			},
			wantReadBack: nil,
			wantErr:      []error{ErrFolderWrite, ErrPathWrite, ErrVaultWrite},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, "") {
				prefix := prefix
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					writeMap := make(map[string]map[string]interface{}, len(tt.give))
					for path, data := range tt.give {
						writeMap[PathJoin(prefix, path)] = data
					}

					err := sharedVaku.FolderWrite(context.Background(), writeMap)
					compareErrors(t, err, tt.wantErr)

					for path, data := range tt.wantReadBack {
						readBack, err := sharedVakuClean.PathRead(PathJoin(prefix, path))
						assert.NoError(t, err)
						assert.Equal(t, data, readBack)
					}
				})
			}
		})
	}
}

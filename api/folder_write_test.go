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
		give         map[string]map[string]any
		wantReadBack map[string]map[string]any
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
			give: map[string]map[string]any{
				"000/001": nil,
			},
			wantReadBack: nil,
			wantErr:      []error{ErrFolderWrite, ErrPathWrite, ErrNilData},
		},
		{
			name: "overwrite",
			give: map[string]map[string]any{
				"0/1": {
					"0001": "0002",
				},
			},
			wantReadBack: map[string]map[string]any{
				"0/1": {
					"0001": "0002",
				},
			},
			wantErr: nil,
		},
		{
			name: "two new paths",
			give: map[string]map[string]any{
				"000/001": {
					"0001": "0002",
				},
				"000/001/002": {
					"0003": "0004",
					"0005": "0006",
				},
			},
			wantReadBack: map[string]map[string]any{
				"000/001": {
					"0001": "0002",
				},
				"000/001/002": {
					"0003": "0004",
					"0005": "0006",
				},
			},
			wantErr: nil,
		},
		{
			name: "two different paths",
			give: map[string]map[string]any{
				"000/001": {
					"0001": "0002",
				},
				"111/001/002": {
					"0003": "0004",
					"0005": "0006",
				},
			},
			wantReadBack: map[string]map[string]any{
				"000/001": {
					"0001": "0002",
				},
				"111/001/002": {
					"0003": "0004",
					"0005": "0006",
				},
			},
			wantErr: nil,
		},
		{
			name: "path write fail",
			give: map[string]map[string]any{
				"failonwrite/error/write/inject": {
					"01": "02",
				},
			},
			wantReadBack: nil,
			wantErr:      []error{ErrFolderWrite, ErrPathWrite, ErrVaultWrite},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, "") {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					writeMap := make(map[string]map[string]any, len(tt.give))
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

package vaku

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderDeleteMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give         string
		wantReadBack map[string]map[string]interface{}
		wantErr      []error
	}{
		{
			give:         "0",
			wantReadBack: nil,
			wantErr:      nil,
		},
		// {
		// 	give:         "0/1",
		// 	wantReadBack: nil,
		// 	wantErr:      nil,
		// },
		// {
		// 	give:         "fake",
		// 	wantReadBack: nil,
		// 	wantErr:      nil,
		// },
		// {
		// 	give: "0/4/13/24/25/error/list/inject",
		// 	wantReadBack: map[string]map[string]interface{}{
		// 		"26/27": {
		// 			"28": "29",
		// 		},
		// 	},
		// 	wantErr: []error{ErrFolderDeleteMeta, ErrFolderListChan, ErrPathList, ErrVaultList},
		// },
		// {
		// 	give: "0/4/13/24/25/26/error/delete/inject",
		// 	wantReadBack: map[string]map[string]interface{}{
		// 		"27": {
		// 			"28": "29",
		// 		},
		// 	},
		// 	wantErr: []error{ErrFolderDeleteMeta, ErrPathDelete, ErrVaultDelete},
		// },
	}

	for _, tt := range tests {
		tt := tt
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				prefix := prefix
				if strings.HasPrefix(prefix, "kv1") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						err := sharedVaku.FolderDeleteMeta(context.Background(), PathJoin(prefix, tt.give))
						compareErrors(t, err, []error{ErrFolderDeleteMeta, ErrPathDeleteMeta, ErrMountVersion})
					})
				}
				if strings.HasPrefix(prefix, "kv2") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						err := sharedVaku.FolderDeleteMeta(context.Background(), PathJoin(prefix, tt.give))
						compareErrors(t, err, tt.wantErr)

						readBack, err := sharedVakuClean.FolderRead(context.Background(), PathJoin(prefix, tt.give))
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReadBack, readBack)
					})
				}
			}
		})
	}
}

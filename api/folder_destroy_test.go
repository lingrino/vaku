package vaku

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderDestroy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give         string
		wantReadBack map[string]map[string]any
		giveVersions []int
		wantErrKV1   []error
		wantErrKV2   []error
	}{
		{
			give:         "0",
			giveVersions: []int{1, 2, 3},
			wantReadBack: nil,
			wantErrKV1:   []error{ErrFolderDestroy, ErrPathDestroy, ErrMountVersion},
			wantErrKV2:   nil,
		},
		{
			give:         "0/1",
			giveVersions: []int{3},
			wantReadBack: nil,
			wantErrKV1:   nil,
			wantErrKV2:   nil,
		},
		{
			give:         "0/4/13/24/25/error/list/inject",
			giveVersions: []int{1, 2},
			wantReadBack: map[string]map[string]any{
				"26/27": {
					"28": "29",
				},
			},
			wantErrKV1: []error{ErrFolderDestroy, ErrFolderListChan, ErrPathList, ErrVaultList},
			wantErrKV2: []error{ErrFolderDestroy, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
		{
			give:         "0/4/13/24/25/26/error/write/inject",
			giveVersions: []int{1, 2},
			wantReadBack: map[string]map[string]any{
				"27": {
					"28": "29",
				},
			},
			wantErrKV1: []error{ErrFolderDestroy, ErrPathDestroy, ErrMountVersion},
			wantErrKV2: []error{ErrFolderDestroy, ErrPathDestroy, ErrVaultWrite},
		},
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

						err := sharedVaku.FolderDestroy(context.Background(), PathJoin(prefix, tt.give), tt.giveVersions)
						compareErrors(t, err, tt.wantErrKV1)
					})
				}
				if strings.HasPrefix(prefix, "kv2") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						err := sharedVaku.FolderDestroy(context.Background(), PathJoin(prefix, tt.give), tt.giveVersions)
						compareErrors(t, err, tt.wantErrKV2)

						readBack, err := sharedVakuClean.FolderRead(context.Background(), PathJoin(prefix, tt.give))
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReadBack, readBack)
					})
				}
			}
		})
	}
}

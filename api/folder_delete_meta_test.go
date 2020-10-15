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
		wantErrKV1   []error
		wantErrKV2   []error
	}{
		{
			give:         "0",
			wantReadBack: nil,
			wantErrKV1:   []error{ErrFolderDeleteMeta, ErrPathDeleteMeta, ErrMountVersion},
			wantErrKV2:   nil,
		},
		{
			give:         "0/1",
			wantReadBack: nil,
			wantErrKV1:   nil,
			wantErrKV2:   nil,
		},
		{
			give:         "fake",
			wantReadBack: nil,
			wantErrKV1:   nil,
			wantErrKV2:   nil,
		},
		{
			give: "0/4/13/24/25/error/list/inject",
			wantReadBack: map[string]map[string]interface{}{
				"26/27": {
					"28": "29",
				},
			},
			wantErrKV1: []error{ErrFolderDeleteMeta, ErrFolderListChan, ErrPathList, ErrVaultList},
			wantErrKV2: []error{ErrFolderDeleteMeta, ErrFolderListChan, ErrPathList, ErrVaultList},
		},
		{
			give: "0/4/13/24/25/26/error/delete/inject",
			wantReadBack: map[string]map[string]interface{}{
				"27": {
					"28": "29",
				},
			},
			wantErrKV1: []error{ErrFolderDeleteMeta, ErrPathDeleteMeta, ErrMountVersion},
			wantErrKV2: []error{ErrFolderDeleteMeta, ErrPathDeleteMeta, ErrVaultDelete},
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

						err := sharedVaku.FolderDeleteMeta(context.Background(), PathJoin(prefix, tt.give))
						compareErrors(t, err, tt.wantErrKV1)
					})
				}
				if strings.HasPrefix(prefix, "kv2") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						err := sharedVaku.FolderDeleteMeta(context.Background(), PathJoin(prefix, tt.give))
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

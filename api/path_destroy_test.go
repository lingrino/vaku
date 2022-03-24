package vaku

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathDestroy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give            string
		giveVersions    []int
		wantErr         []error
		wantNilReadback bool
	}{
		{
			give:         "0/1",
			giveVersions: nil,
			wantErr:      []error{ErrPathDestroy},
		},
		{
			give:         "0/1",
			giveVersions: []int{},
			wantErr:      []error{ErrPathDestroy},
		},
		{
			give:         "0/1",
			giveVersions: []int{1},
		},
		{
			give:            "0/1",
			giveVersions:    []int{2},
			wantNilReadback: true,
		},
		{
			give:            "fake",
			wantErr:         nil,
			giveVersions:    []int{1, 2, 3, 4, 5, 6, 7},
			wantNilReadback: true,
		},
		{
			give:         "error/write/inject",
			giveVersions: []int{1},
			wantErr:      []error{ErrPathDestroy, ErrVaultWrite},
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

						err := sharedVaku.PathDestroy(PathJoin(prefix, tt.give), []int{1})
						compareErrors(t, err, []error{ErrPathDestroy, ErrMountVersion})
					})
				}
				if strings.HasPrefix(prefix, "kv2") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						// overwrite all paths to create a new version
						overwrite := map[string]any{"foo": "bar"}
						err := sharedVakuClean.PathWrite(PathJoin(prefix, tt.give), overwrite)
						assert.NoError(t, err)

						err = sharedVaku.PathDestroy(PathJoin(prefix, tt.give), tt.giveVersions)
						compareErrors(t, err, tt.wantErr)

						readBack, err := sharedVakuClean.PathRead(PathJoin(prefix, tt.give))
						assert.NoError(t, err)
						if tt.wantNilReadback {
							assert.Nil(t, readBack)
						} else {
							assert.Equal(t, overwrite, readBack)
						}
					})
				}
			}
		})
	}
}

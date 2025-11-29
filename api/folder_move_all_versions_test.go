package vaku

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderMoveAllVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveSrc    string
		giveDst    string
		wantErr    []error
		wantNilSrc bool
		wantNilDst bool
	}{
		{
			giveSrc:    "0/1",
			giveDst:    "moveallversions/0/1",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			giveSrc:    "0",
			giveDst:    "moveallversions/0",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			giveSrc: "0/4/13/24/25/26/error/list/inject",
			giveDst: "moveallversions/error/list",
			wantErr: []error{
				ErrFolderMoveAllVersions, ErrFolderCopyAllVersions,
				ErrFolderListChan, ErrPathList, ErrVaultList,
			},
			wantNilDst: true,
		},
		{
			giveSrc: "0/4/13/24/25/26/error/delete/inject",
			giveDst: "moveallversions/error/delete",
			wantErr: []error{ErrFolderMoveAllVersions, ErrFolderDeleteMeta, ErrPathDeleteMeta, ErrVaultDelete},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				// Skip kv1 mounts - FolderMoveAllVersions only works on kv2
				if strings.HasPrefix(prefixPair[0], "kv1") || strings.HasPrefix(prefixPair[1], "kv1") {
					continue
				}

				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					pathSrc := PathJoin(prefixPair[0], tt.giveSrc)
					pathDst := PathJoin(prefixPair[1], tt.giveDst)

					origSrc, err := sharedVakuClean.FolderRead(context.Background(), pathSrc)
					assert.NoError(t, err)
					TrimPrefixMap(origSrc, pathSrc)

					err = sharedVaku.FolderMoveAllVersions(context.Background(), pathSrc, pathDst)
					compareErrors(t, err, tt.wantErr)

					readSrc, errSrc := sharedVakuClean.FolderRead(context.Background(), pathSrc)
					readDst, errDst := sharedVakuClean.dc.FolderRead(context.Background(), pathDst)
					assert.NoError(t, errSrc)
					assert.NoError(t, errDst)
					TrimPrefixMap(readSrc, pathSrc)
					TrimPrefixMap(readDst, pathDst)

					if tt.wantNilSrc {
						assert.Nil(t, readSrc)
					} else {
						assert.Equal(t, origSrc, readSrc)
					}
					if tt.wantNilDst {
						assert.Nil(t, readDst)
					} else {
						assert.Equal(t, origSrc, readDst)
					}
				})
			}
		})
	}
}

func TestFolderMoveAllVersionsKV1(t *testing.T) {
	t.Parallel()

	// Test that KV1 source returns error
	t.Run("kv1 source", func(t *testing.T) {
		t.Parallel()
		for _, prefix := range seededPrefixes(t, "0/1") {
			if strings.HasPrefix(prefix, "kv1") {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.FolderMoveAllVersions(
						context.Background(),
						PathJoin(prefix, "0/1"),
						"kv2/moveallversions/dst",
					)
					compareErrors(t, err, []error{
						ErrFolderMoveAllVersions,
						ErrFolderCopyAllVersions,
						ErrMountVersion,
					})
				})
			}
		}
	})

	// Test that KV1 destination returns error
	t.Run("kv1 destination", func(t *testing.T) {
		t.Parallel()
		for _, prefix := range seededPrefixes(t, "0/1") {
			if strings.HasPrefix(prefix, "kv2") {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.FolderMoveAllVersions(
						context.Background(),
						PathJoin(prefix, "0/1"),
						"kv1/moveallversions/dst",
					)
					compareErrors(t, err, []error{
						ErrFolderMoveAllVersions,
						ErrFolderCopyAllVersions,
						ErrMountVersion,
					})
				})
			}
		}
	})
}

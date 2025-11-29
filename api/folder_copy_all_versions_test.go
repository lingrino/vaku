package vaku

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderCopyAllVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveSrc    string
		giveDst    string
		wantErr    []error
		wantNilDst bool
	}{
		{
			giveSrc: "0/1",
			giveDst: "copyallversions/0/1",
			wantErr: nil,
		},
		{
			giveSrc: "0",
			giveDst: "copyallversions/0",
			wantErr: nil,
		},
		{
			giveSrc:    "0/4/13/24/25/26/error/list/inject",
			giveDst:    "copyallversions/error/list",
			wantErr:    []error{ErrFolderCopyAllVersions, ErrFolderListChan, ErrPathList, ErrVaultList},
			wantNilDst: true,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				// Skip kv1 mounts - FolderCopyAllVersions only works on kv2
				if strings.HasPrefix(prefixPair[0], "kv1") || strings.HasPrefix(prefixPair[1], "kv1") {
					continue
				}

				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					pathSrc := PathJoin(prefixPair[0], tt.giveSrc)
					pathDst := PathJoin(prefixPair[1], tt.giveDst)

					err := sharedVaku.FolderCopyAllVersions(context.Background(), pathSrc, pathDst)
					compareErrors(t, err, tt.wantErr)

					if tt.wantNilDst {
						readDst, errDst := sharedVakuClean.dc.FolderRead(context.Background(), pathDst)
						assert.NoError(t, errDst)
						assert.Nil(t, readDst)
					} else if tt.wantErr == nil {
						// Verify source and destination have matching current data
						readSrc, errSrc := sharedVakuClean.FolderRead(context.Background(), pathSrc)
						readDst, errDst := sharedVakuClean.dc.FolderRead(context.Background(), pathDst)
						assert.NoError(t, errSrc)
						assert.NoError(t, errDst)

						TrimPrefixMap(readSrc, prefixPair[0])
						TrimPrefixMap(readDst, prefixPair[1])
						assert.Equal(t, readSrc, readDst)
					}
				})
			}
		})
	}
}

func TestFolderCopyAllVersionsKV1(t *testing.T) {
	t.Parallel()

	// Test that KV1 source returns error
	t.Run("kv1 source", func(t *testing.T) {
		t.Parallel()
		for _, prefix := range seededPrefixes(t, "0/1") {
			if strings.HasPrefix(prefix, "kv1") {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.FolderCopyAllVersions(
						context.Background(),
						PathJoin(prefix, "0/1"),
						"kv2/copyallversions/dst",
					)
					compareErrors(t, err, []error{
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

					err := sharedVaku.FolderCopyAllVersions(
						context.Background(),
						PathJoin(prefix, "0/1"),
						"kv1/copyallversions/dst",
					)
					compareErrors(t, err, []error{
						ErrFolderCopyAllVersions,
						ErrMountVersion,
					})
				})
			}
		}
	})
}

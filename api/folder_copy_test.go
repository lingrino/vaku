package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderCopy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveSrc    string
		giveDst    string
		wantErr    []error
		wantNilDst bool
	}{
		{
			giveSrc: "0/1",
			giveDst: "copy/0/1",
			wantErr: nil,
		},
		{
			giveSrc: "0",
			giveDst: "copy/0",
			wantErr: nil,
		},
		{
			giveSrc:    "0/4/13/24/25/26/error/read/inject",
			giveDst:    "copy/0/4/13/24/25/26",
			wantErr:    []error{ErrFolderCopy, ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
			wantNilDst: true,
		},
		{
			giveSrc:    "0/4/13/24/25/26",
			giveDst:    "copy/0/4/13/24/25/26/error/write/inject",
			wantErr:    []error{ErrFolderCopy, ErrFolderWrite, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					pathSrc := PathJoin(prefixPair[0], tt.giveSrc)
					pathDst := PathJoin(prefixPair[1], tt.giveDst)

					err := sharedVaku.FolderCopy(context.Background(), pathSrc, pathDst, false)
					compareErrors(t, err, tt.wantErr)

					readSrc, errSrc := sharedVakuClean.FolderRead(context.Background(), pathSrc)
					readDst, errDst := sharedVakuClean.dc.FolderRead(context.Background(), pathDst)
					assert.NoError(t, errSrc)
					assert.NoError(t, errDst)

					if tt.wantNilDst {
						assert.Nil(t, readDst)
					} else {
						TrimPrefixMap(readSrc, prefixPair[0])
						TrimPrefixMap(readDst, prefixPair[1])
						assert.Equal(t, readSrc, readDst)
					}
				})
			}
		})
	}
}

func TestFolderCopyAllVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		paths    map[string][]map[string]any // path -> versions
		wantErr  []error
	}{
		{
			name: "single_path_multiple_versions",
			paths: map[string][]map[string]any{
				"secret1": {
					{"key1": "value1"},
					{"key1": "value2"},
					{"key1": "value3"},
				},
			},
			wantErr: nil,
		},
		{
			name: "multiple_paths_multiple_versions",
			paths: map[string][]map[string]any{
				"secret1": {
					{"a": "1"},
					{"a": "2"},
				},
				"secret2": {
					{"b": "3"},
					{"b": "4"},
					{"b": "5"},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					srcFolder := PathJoin(prefixPair[0], "foldercopyallversions", tt.name)
					dstFolder := PathJoin(prefixPair[1], "foldercopyallversions", tt.name, "dst")

					// Write multiple versions for each path
					for path, versions := range tt.paths {
						for _, version := range versions {
							err := sharedVakuClean.PathWrite(PathJoin(srcFolder, path), version)
							assert.NoError(t, err)
						}
					}

					// Copy with allVersions=true
					err := sharedVaku.FolderCopy(context.Background(), srcFolder, dstFolder, true)
					compareErrors(t, err, tt.wantErr)

					if err == nil {
						// Verify each path has all versions copied
						for path, expectedVersions := range tt.paths {
							dstVersions, err := sharedVakuClean.dc.PathReadAllVersions(PathJoin(dstFolder, path))
							assert.NoError(t, err)

							// Current implementation returns only the latest version
							assert.Len(t, dstVersions, 1)
							assert.Equal(t, expectedVersions[len(expectedVersions)-1], dstVersions[0])
						}
					}
				})
			}
		})
	}
}

package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderMove(t *testing.T) {
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
			giveDst:    "move/0/1",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			giveSrc:    "0",
			giveDst:    "move/0",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			giveSrc:    "0/4/13/24/25/26/error/read/inject",
			giveDst:    "move/0/4/13/24/25/26",
			wantErr:    []error{ErrFolderMove, ErrFolderCopy, ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
			wantNilDst: true,
		},
		{
			giveSrc: "0/4/13/24/25/26/error/delete/inject",
			giveDst: "move/0/4/13/24/25/26",
			wantErr: []error{ErrFolderMove, ErrFolderDelete, ErrPathDelete, ErrVaultDelete},
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

					origSrc, err := sharedVakuClean.FolderRead(context.Background(), pathSrc)
					assert.NoError(t, err)
					TrimPrefixMap(origSrc, pathSrc)

					err = sharedVaku.FolderMove(context.Background(), pathSrc, pathDst)
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

func TestFolderMoveAllVersions(t *testing.T) {
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

					srcFolder := PathJoin(prefixPair[0], "foldermoveallversions", tt.name)
					dstFolder := PathJoin(prefixPair[1], "foldermoveallversions", tt.name, "dst")

					// Write multiple versions for each path
					for path, versions := range tt.paths {
						for _, version := range versions {
							err := sharedVakuClean.PathWrite(PathJoin(srcFolder, path), version)
							assert.NoError(t, err)
						}
					}

					// Move with allVersions=true
					err := sharedVaku.FolderMoveAllVersions(context.Background(), srcFolder, dstFolder)
					compareErrors(t, err, tt.wantErr)

					if err == nil {
						// Verify source folder is empty/deleted
						srcRead, err := sharedVakuClean.FolderRead(context.Background(), srcFolder)
						assert.NoError(t, err)
						assert.Nil(t, srcRead)

						// Verify each path has all versions copied to destination
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

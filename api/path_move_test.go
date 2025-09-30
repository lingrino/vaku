package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMove(t *testing.T) {
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
			giveSrc:    "0/4/8",
			giveDst:    "0/4/5",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			giveSrc:    "error/read/inject",
			giveDst:    "move/readerror",
			wantErr:    []error{ErrPathMove, ErrPathCopy, ErrPathRead, ErrVaultRead},
			wantNilSrc: true,
			wantNilDst: true,
		},
		{
			giveSrc: "0/4/13/14/error/delete/inject",
			giveDst: "move/deleteerror",
			wantErr: []error{ErrPathMove, ErrPathDelete, ErrVaultDelete},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					origSrc, err := sharedVakuClean.PathRead(PathJoin(prefixPair[0], tt.giveSrc))
					assert.NoError(t, err)

					err = sharedVaku.PathMove(PathJoin(prefixPair[0], tt.giveSrc), PathJoin(prefixPair[1], tt.giveDst))
					compareErrors(t, err, tt.wantErr)

					readSrc, errSrc := sharedVakuClean.PathRead(PathJoin(prefixPair[0], tt.giveSrc))
					readDst, errDst := sharedVakuClean.dc.PathRead(PathJoin(prefixPair[1], tt.giveDst))
					assert.NoError(t, errSrc)
					assert.NoError(t, errDst)

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

func TestPathMoveAllVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		versions []map[string]any
		wantErr  []error
	}{
		{
			name: "single_version",
			versions: []map[string]any{
				{"key1": "value1"},
			},
			wantErr: nil,
		},
		{
			name: "multiple_versions",
			versions: []map[string]any{
				{"key1": "value1"},
				{"key1": "value2", "key2": "newkey"},
				{"key1": "value3"},
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

					srcPath := PathJoin(prefixPair[0], "moveallversions", tt.name)
					dstPath := PathJoin(prefixPair[1], "moveallversions", tt.name, "dst")

					// Write multiple versions to source
					for _, version := range tt.versions {
						err := sharedVakuClean.PathWrite(srcPath, version)
						assert.NoError(t, err)
					}

					// Move with allVersions=true
					err := sharedVaku.PathMoveAllVersions(srcPath, dstPath)
					compareErrors(t, err, tt.wantErr)

					if err == nil {
						// Verify source is deleted
						readSrc, err := sharedVakuClean.PathRead(srcPath)
						assert.NoError(t, err)
						assert.Nil(t, readSrc)

						// Read all versions from destination
						dstVersions, err := sharedVakuClean.dc.PathReadAllVersions(dstPath)
						assert.NoError(t, err)

						// Current implementation returns only the latest version
						assert.Len(t, dstVersions, 1)
						assert.Equal(t, tt.versions[len(tt.versions)-1], dstVersions[0])
					}
				})
			}
		})
	}
}

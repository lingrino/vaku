package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathCopy(t *testing.T) {
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
			giveSrc: "0/1",
			giveDst: "0/4/5",
			wantErr: nil,
		},
		{
			giveSrc:    "0/4/8/error/read/inject",
			giveDst:    "copy/readerror",
			wantErr:    []error{ErrPathCopy, ErrPathRead, ErrVaultRead},
			wantNilDst: true,
		},
		{
			giveSrc:    "0/4/8",
			giveDst:    "copy/writeerror/error/write/inject",
			wantErr:    []error{ErrPathCopy, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.PathCopy(PathJoin(prefixPair[0], tt.giveSrc), PathJoin(prefixPair[1], tt.giveDst), false)
					compareErrors(t, err, tt.wantErr)

					readSrc, errSrc := sharedVakuClean.PathRead(PathJoin(prefixPair[0], tt.giveSrc))
					readDst, errDst := sharedVakuClean.dc.PathRead(PathJoin(prefixPair[1], tt.giveDst))
					assert.NoError(t, errSrc)
					assert.NoError(t, errDst)

					if tt.wantNilDst {
						assert.Nil(t, readDst)
					} else {
						assert.Equal(t, readSrc, readDst)
					}
				})
			}
		})
	}
}

func TestPathCopyAllVersions(t *testing.T) {
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

					srcPath := PathJoin(prefixPair[0], "copyallversions", tt.name)
					dstPath := PathJoin(prefixPair[1], "copyallversions", tt.name, "dst")

					// Write multiple versions to source
					for _, version := range tt.versions {
						err := sharedVakuClean.PathWrite(srcPath, version)
						assert.NoError(t, err)
					}

					// Copy with allVersions=true
					err := sharedVaku.PathCopy(srcPath, dstPath, true)
					compareErrors(t, err, tt.wantErr)

					if err == nil {
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

package vaku

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathCopyAllVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		giveSrc    string
		giveDst    string
		wantErr    []error
		wantNilDst bool
	}{
		{
			name:    "basic copy",
			giveSrc: "0/1",
			giveDst: "copy/allversions/0/1",
			wantErr: nil,
		},
		{
			name:    "copy to different path",
			giveSrc: "0/4/5",
			giveDst: "copy/allversions/different",
			wantErr: nil,
		},
		{
			name:       "source not found",
			giveSrc:    "fake/nonexistent",
			giveDst:    "copy/allversions/fake",
			wantErr:    nil,
			wantNilDst: true,
		},
		{
			name:       "read error",
			giveSrc:    "0/4/8/error/read/inject",
			giveDst:    "copy/allversions/readerror",
			wantErr:    []error{ErrPathCopyAllVersions, ErrPathReadMeta, ErrVaultRead},
			wantNilDst: true,
		},
		{
			name:       "write error",
			giveSrc:    "0/4/8",
			giveDst:    "copy/allversions/error/write/inject",
			wantErr:    []error{ErrPathCopyAllVersions, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				// Skip kv1 mounts - PathCopyAllVersions only works on kv2
				if strings.HasPrefix(prefixPair[0], "kv1") || strings.HasPrefix(prefixPair[1], "kv1") {
					continue
				}

				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					srcPath := PathJoin(prefixPair[0], tt.giveSrc)
					dstPath := PathJoin(prefixPair[1], tt.giveDst)

					err := sharedVaku.PathCopyAllVersions(srcPath, dstPath)
					compareErrors(t, err, tt.wantErr)

					if tt.wantNilDst {
						readDst, errDst := sharedVakuClean.dc.PathRead(dstPath)
						assert.NoError(t, errDst)
						assert.Nil(t, readDst)
					} else if tt.wantErr == nil {
						// Verify source and destination have matching current data
						readSrc, errSrc := sharedVakuClean.PathRead(srcPath)
						readDst, errDst := sharedVakuClean.dc.PathRead(dstPath)
						assert.NoError(t, errSrc)
						assert.NoError(t, errDst)
						assert.Equal(t, readSrc, readDst)

						// Verify version counts match
						srcMeta, errSrcMeta := sharedVakuClean.PathReadMeta(srcPath)
						dstMeta, errDstMeta := sharedVakuClean.dc.PathReadMeta(dstPath)
						assert.NoError(t, errSrcMeta)
						assert.NoError(t, errDstMeta)
						if srcMeta != nil && dstMeta != nil {
							assert.Equal(t, len(srcMeta.Versions), len(dstMeta.Versions))
						}
					}
				})
			}
		})
	}
}

func TestPathCopyAllVersionsKV1(t *testing.T) {
	t.Parallel()

	// Test that KV1 source returns error
	t.Run("kv1 source", func(t *testing.T) {
		t.Parallel()
		for _, prefix := range seededPrefixes(t, "0/1") {
			if strings.HasPrefix(prefix, "kv1") {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.PathCopyAllVersions(PathJoin(prefix, "0/1"), "kv2/copy/dst")
					compareErrors(t, err, []error{ErrPathCopyAllVersions, ErrMountVersion})
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

					err := sharedVaku.PathCopyAllVersions(PathJoin(prefix, "0/1"), "kv1/copy/dst")
					compareErrors(t, err, []error{ErrPathCopyAllVersions, ErrMountVersion})
				})
			}
		}
	})
}

func TestPathCopyAllVersionsWithDeletedVersions(t *testing.T) {
	t.Parallel()

	for _, prefix := range seededPrefixes(t, "multiversion") {
		if !strings.HasPrefix(prefix, "kv2") {
			continue
		}

		t.Run(testName(prefix), func(t *testing.T) {
			t.Parallel()

			srcPath := PathJoin(prefix, "multiversion/test")
			dstPath := PathJoin(prefix, "copy/multiversion/test")

			// Create multiple versions
			version1 := map[string]any{"version": "1", "data": "first"}
			version2 := map[string]any{"version": "2", "data": "second"}
			version3 := map[string]any{"version": "3", "data": "third"}

			err := sharedVakuClean.PathWrite(srcPath, version1)
			assert.NoError(t, err)
			err = sharedVakuClean.PathWrite(srcPath, version2)
			assert.NoError(t, err)
			err = sharedVakuClean.PathWrite(srcPath, version3)
			assert.NoError(t, err)

			// Verify we have 3 versions
			srcMeta, err := sharedVakuClean.PathReadMeta(srcPath)
			assert.NoError(t, err)
			if assert.NotNil(t, srcMeta) {
				assert.Equal(t, 3, len(srcMeta.Versions))
			}

			// Copy all versions
			err = sharedVaku.PathCopyAllVersions(srcPath, dstPath)
			assert.NoError(t, err)

			// Verify destination has 3 versions (use dc since dst is on destination Vault)
			dstMeta, err := sharedVakuClean.dc.PathReadMeta(dstPath)
			assert.NoError(t, err)
			if assert.NotNil(t, dstMeta) {
				assert.Equal(t, 3, len(dstMeta.Versions))
			}

			// Verify each version's data by reading specific versions
			dstV1, err := sharedVakuClean.dc.PathReadVersion(dstPath, 1)
			assert.NoError(t, err)
			assert.Equal(t, version1, dstV1)

			dstV2, err := sharedVakuClean.dc.PathReadVersion(dstPath, 2)
			assert.NoError(t, err)
			assert.Equal(t, version2, dstV2)

			dstV3, err := sharedVakuClean.dc.PathReadVersion(dstPath, 3)
			assert.NoError(t, err)
			assert.Equal(t, version3, dstV3)
		})
	}
}

func TestPathCopyAllVersionsWithDestroyedVersion(t *testing.T) {
	t.Parallel()

	for _, prefix := range seededPrefixes(t, "destroyed") {
		if !strings.HasPrefix(prefix, "kv2") {
			continue
		}

		t.Run(testName(prefix), func(t *testing.T) {
			t.Parallel()

			srcPath := PathJoin(prefix, "destroyed/test")
			dstPath := PathJoin(prefix, "copy/destroyed/test")

			// Create multiple versions
			version1 := map[string]any{"version": "1"}
			version2 := map[string]any{"version": "2"}
			version3 := map[string]any{"version": "3"}

			err := sharedVakuClean.PathWrite(srcPath, version1)
			assert.NoError(t, err)
			err = sharedVakuClean.PathWrite(srcPath, version2)
			assert.NoError(t, err)
			err = sharedVakuClean.PathWrite(srcPath, version3)
			assert.NoError(t, err)

			// Destroy version 2
			err = sharedVakuClean.PathDestroy(srcPath, []int{2})
			assert.NoError(t, err)

			// Verify version 2 is destroyed
			srcMeta, err := sharedVakuClean.PathReadMeta(srcPath)
			assert.NoError(t, err)
			if assert.NotNil(t, srcMeta) {
				assert.True(t, srcMeta.Versions[2].Destroyed)
			}

			// Copy all versions
			err = sharedVaku.PathCopyAllVersions(srcPath, dstPath)
			assert.NoError(t, err)

			// Verify destination has 3 versions (use dc since dst is on destination Vault)
			dstMeta, err := sharedVakuClean.dc.PathReadMeta(dstPath)
			assert.NoError(t, err)
			if assert.NotNil(t, dstMeta) {
				assert.Equal(t, 3, len(dstMeta.Versions))
			}

			// Verify version 1 data
			dstV1, err := sharedVakuClean.dc.PathReadVersion(dstPath, 1)
			assert.NoError(t, err)
			assert.Equal(t, version1, dstV1)

			// Verify version 2 is empty (was destroyed at source)
			dstV2, err := sharedVakuClean.dc.PathReadVersion(dstPath, 2)
			assert.NoError(t, err)
			assert.Empty(t, dstV2)

			// Verify version 3 data
			dstV3, err := sharedVakuClean.dc.PathReadVersion(dstPath, 3)
			assert.NoError(t, err)
			assert.Equal(t, version3, dstV3)
		})
	}
}

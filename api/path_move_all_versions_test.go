package vaku

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMoveAllVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		giveSrc    string
		giveDst    string
		wantErr    []error
		wantNilSrc bool
		wantNilDst bool
	}{
		{
			name:       "basic move",
			giveSrc:    "0/1",
			giveDst:    "move/allversions/0/1",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			name:       "move to different path",
			giveSrc:    "0/4/5",
			giveDst:    "move/allversions/different",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			name:       "source not found",
			giveSrc:    "fake/nonexistent",
			giveDst:    "move/allversions/fake",
			wantErr:    nil,
			wantNilSrc: true,
			wantNilDst: true,
		},
		{
			name:       "read error",
			giveSrc:    "error/read/inject",
			giveDst:    "move/allversions/readerror",
			wantErr:    []error{ErrPathMoveAllVersions, ErrPathCopyAllVersions, ErrPathReadMeta, ErrVaultRead},
			wantNilSrc: true,
			wantNilDst: true,
		},
		{
			name:       "write error",
			giveSrc:    "0/4/8",
			giveDst:    "move/allversions/error/write/inject",
			wantErr:    []error{ErrPathMoveAllVersions, ErrPathCopyAllVersions, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
		{
			name:       "delete error",
			giveSrc:    "0/4/13/14/error/delete/inject",
			giveDst:    "move/allversions/deleteerror",
			wantErr:    []error{ErrPathMoveAllVersions, ErrPathDeleteMeta, ErrVaultDelete},
			wantNilSrc: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				// Skip kv1 mounts - PathMoveAllVersions only works on kv2
				if strings.HasPrefix(prefixPair[0], "kv1") || strings.HasPrefix(prefixPair[1], "kv1") {
					continue
				}

				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					srcPath := PathJoin(prefixPair[0], tt.giveSrc)
					dstPath := PathJoin(prefixPair[1], tt.giveDst)

					origSrc, err := sharedVakuClean.PathRead(srcPath)
					assert.NoError(t, err)

					err = sharedVaku.PathMoveAllVersions(srcPath, dstPath)
					compareErrors(t, err, tt.wantErr)

					readSrc, errSrc := sharedVakuClean.PathRead(srcPath)
					readDst, errDst := sharedVakuClean.dc.PathRead(dstPath)
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

func TestPathMoveAllVersionsKV1(t *testing.T) {
	t.Parallel()

	// Test that KV1 source returns error
	t.Run("kv1 source", func(t *testing.T) {
		t.Parallel()
		for _, prefix := range seededPrefixes(t, "0/1") {
			if strings.HasPrefix(prefix, "kv1") {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.PathMoveAllVersions(PathJoin(prefix, "0/1"), "kv2/move/dst")
					compareErrors(t, err, []error{ErrPathMoveAllVersions, ErrPathCopyAllVersions, ErrMountVersion})
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

					err := sharedVaku.PathMoveAllVersions(PathJoin(prefix, "0/1"), "kv1/move/dst")
					compareErrors(t, err, []error{ErrPathMoveAllVersions, ErrPathCopyAllVersions, ErrMountVersion})
				})
			}
		}
	})
}

func TestPathMoveAllVersionsWithMultipleVersions(t *testing.T) {
	t.Parallel()

	for _, prefix := range seededPrefixes(t, "movemulti") {
		if !strings.HasPrefix(prefix, "kv2") {
			continue
		}

		t.Run(testName(prefix), func(t *testing.T) {
			t.Parallel()

			srcPath := PathJoin(prefix, "movemulti/test")
			dstPath := PathJoin(prefix, "move/movemulti/test")

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

			// Move all versions
			err = sharedVaku.PathMoveAllVersions(srcPath, dstPath)
			assert.NoError(t, err)

			// Verify source is deleted
			srcRead, err := sharedVakuClean.PathRead(srcPath)
			assert.NoError(t, err)
			assert.Nil(t, srcRead)

			// Verify destination has 3 versions (use dc since dst is on destination Vault)
			dstMeta, err := sharedVakuClean.dc.PathReadMeta(dstPath)
			assert.NoError(t, err)
			if assert.NotNil(t, dstMeta) {
				assert.Equal(t, 3, len(dstMeta.Versions))
			}

			// Verify each version's data
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

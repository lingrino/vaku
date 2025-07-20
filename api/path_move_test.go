package vaku

import (
	"strings"
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
			giveSrc:    "0/4/8",
			giveDst:    "0/4/5",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			giveSrc:    "error/read/inject",
			giveDst:    "moveallversions/readerror",
			wantErr:    []error{ErrPathMoveAllVersions, ErrPathCopyAllVersions, ErrPathReadMetadata, ErrVaultRead},
			wantNilSrc: true,
			wantNilDst: true,
		},
		{
			giveSrc:    "fake",
			giveDst:    "moveallversions/fake",
			wantErr:    nil,
			wantNilSrc: true,
			wantNilDst: true,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
					t.Parallel()

					// Only test on v2 mounts since all versions is v2 only
					if !strings.Contains(prefixPair[0], "v2") || !strings.Contains(prefixPair[1], "v2") {
						t.Skip("PathMoveAllVersions only works on KV v2 mounts")
					}

					origSrc, err := sharedVakuClean.PathRead(PathJoin(prefixPair[0], tt.giveSrc))
					assert.NoError(t, err)

					err = sharedVaku.PathMoveAllVersions(PathJoin(prefixPair[0], tt.giveSrc), PathJoin(prefixPair[1], tt.giveDst))
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

					// For successful moves, verify that source metadata is completely gone
					if err == nil && !tt.wantNilSrc && !tt.wantNilDst {
						srcMetadata, metaErr := sharedVakuClean.PathReadMetadata(PathJoin(prefixPair[0], tt.giveSrc))
						assert.NoError(t, metaErr)
						assert.Nil(t, srcMetadata, "Source metadata should be completely deleted after move")
					}
				})
			}
		})
	}
}

func TestPathMoveAllVersionsKV1(t *testing.T) {
	t.Parallel()

	// Test that PathMoveAllVersions returns proper error on KV v1
	for _, prefixPair := range seededPrefixProduct(t) {
		t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
			t.Parallel()

			// Only test on v1 mounts
			if !strings.Contains(prefixPair[0], "v1") {
				t.Skip("Testing KV v1 error handling")
			}

			srcPath := PathJoin(prefixPair[0], "0/1")
			dstPath := PathJoin(prefixPair[1], "moveallversions/kv1test")
			err := sharedVaku.PathMoveAllVersions(srcPath, dstPath)
			expectedErrors := []error{ErrPathMoveAllVersions, ErrMountVersion}
			compareErrors(t, err, expectedErrors)
		})
	}
}

package vaku

import (
	"strings"
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

					err := sharedVaku.PathCopy(PathJoin(prefixPair[0], tt.giveSrc), PathJoin(prefixPair[1], tt.giveDst))
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
			giveSrc: "0/1",
			giveDst: "0/4/5",
			wantErr: nil,
		},
		{
			giveSrc:    "0/4/8/error/read/inject",
			giveDst:    "copyallversions/readerror",
			wantErr:    []error{ErrPathCopyAllVersions, ErrPathReadMetadata, ErrVaultRead},
			wantNilDst: true,
		},
		{
			giveSrc:    "0/4/8",
			giveDst:    "copyallversions/writeerror/error/write/inject",
			wantErr:    []error{ErrPathCopyAllVersions, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
		{
			giveSrc:    "fake",
			giveDst:    "copyallversions/fake",
			wantErr:    nil,
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
						t.Skip("PathCopyAllVersions only works on KV v2 mounts")
					}

					err := sharedVaku.PathCopyAllVersions(PathJoin(prefixPair[0], tt.giveSrc), PathJoin(prefixPair[1], tt.giveDst))
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

func TestPathCopyAllVersionsKV1(t *testing.T) {
	t.Parallel()

	// Test that PathCopyAllVersions returns proper error on KV v1
	for _, prefixPair := range seededPrefixProduct(t) {
		t.Run(testName(prefixPair[0], prefixPair[1]), func(t *testing.T) {
			t.Parallel()

			// Only test on v1 mounts
			if !strings.Contains(prefixPair[0], "v1") {
				t.Skip("Testing KV v1 error handling")
			}

			srcPath := PathJoin(prefixPair[0], "0/1")
			dstPath := PathJoin(prefixPair[1], "copyallversions/kv1test")
			err := sharedVaku.PathCopyAllVersions(srcPath, dstPath)
			expectedErrors := []error{ErrPathCopyAllVersions, ErrMountVersion}
			compareErrors(t, err, expectedErrors)
		})
	}
}

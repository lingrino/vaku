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
		tt := tt
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				prefixPair := prefixPair
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

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
		tt := tt
		t.Run(testName(tt.giveSrc, tt.giveDst), func(t *testing.T) {
			t.Parallel()
			for _, prefixPair := range seededPrefixProduct(t) {
				prefixPair := prefixPair
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

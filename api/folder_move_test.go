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

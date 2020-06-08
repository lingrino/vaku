package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderCopy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		giveSrc    string
		giveDst    string
		wantErr    []error
		wantNilDst bool
	}{
		// {
		// 	giveSrc: "0/1",
		// 	giveDst: "copy/0/1",
		// 	wantErr: nil,
		// },
		// {
		// 	giveSrc: "0",
		// 	giveDst: "copy/0",
		// 	wantErr: nil,
		// },
		// {
		// 	giveSrc:    "0/4/13/24/25/26/error/read/inject",
		// 	giveDst:    "copy/0/4/13/24/25/26",
		// 	wantErr:    []error{ErrFolderCopy, ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
		// 	wantNilDst: true,
		// },
		{
			giveSrc:    "0/4/13/24/25/26",
			giveDst:    "copy/0/4/13/24/25/26/error/write/inject",
			wantErr:    []error{ErrFolderCopy, ErrFolderWrite, ErrPathWrite, ErrVaultWrite},
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

					pathSrc := PathJoin(prefixPair[0], tt.giveSrc)
					pathDst := PathJoin(prefixPair[1], tt.giveDst)

					err := sharedVaku.FolderCopy(context.Background(), pathSrc, pathDst)
					compareErrors(t, err, tt.wantErr)

					readSrc, errSrc := sharedVakuClean.FolderRead(context.Background(), pathSrc)
					readDst, errDst := sharedVakuClean.dc.FolderRead(context.Background(), pathDst)
					assert.NoError(t, errSrc)
					assert.NoError(t, errDst)

					if tt.wantNilDst {
						assert.Nil(t, readDst)
					} else {
						TrimPrefixMap(readSrc, prefixPair[0])
						TrimPrefixMap(readDst, prefixPair[1])
						assert.Equal(t, readSrc, readDst)
					}
				})
			}
		})
	}
}

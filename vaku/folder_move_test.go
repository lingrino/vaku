package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderMove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		giveSrc        string
		giveDst        string
		giveSrcLogical logical
		giveDstLogical logical
		giveOptions    []Option
		wantErr        []error
		wantNilSrc     bool
		wantNilDst     bool
	}{
		{
			name:       "move one",
			giveSrc:    "test/foo",
			giveDst:    "moveone/test",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			name:       "move all",
			giveSrc:    "test",
			giveDst:    "moveall/test",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			name:        "move all absolute path",
			giveSrc:     "test",
			giveDst:     "move/test",
			giveOptions: []Option{WithabsolutePath(true)},
			wantErr:     nil,
			wantNilSrc:  true,
		},
		{
			name:    "read err",
			giveSrc: "test/inner/again",
			giveDst: "readerr/test",
			giveSrcLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:    []error{ErrFolderMove, ErrFolderCopy, ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
			wantNilDst: true,
		},
		{
			name:    "write err",
			giveSrc: "test/inner/again",
			giveDst: "writerr/test",
			giveSrcLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantErr:    []error{ErrFolderMove, ErrFolderCopy, ErrFolderWrite, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
		{
			name:    "delete err",
			giveSrc: "test/inner/again",
			giveDst: "writerr/test",
			giveSrcLogical: &errLogical{
				err: errInject,
				op:  "Delete",
			},
			wantErr: []error{ErrFolderMove, ErrFolderDelete, ErrPathDelete, ErrVaultDelete},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, ver := range versionProduct {
				client := testClient(t, tt.giveOptions...)
				clientDD := testClientDiffDst(t, tt.giveOptions...)

				for _, c := range []*Client{client, clientDD} {
					readbackClient := cloneCLient(t, c)
					updateLogical(t, c, tt.giveSrcLogical, tt.giveDstLogical)

					pathS := addMountToPath(t, tt.giveSrc, ver[0])
					pathD := addMountToPath(t, tt.giveDst, ver[1])

					orig, err := readbackClient.FolderRead(context.Background(), pathS)
					assert.NoError(t, err)
					TrimPrefixMap(orig, pathS)

					err = c.FolderMove(context.Background(), pathS, pathD)
					compareErrors(t, err, tt.wantErr)

					readBackS, errS := readbackClient.FolderRead(context.Background(), pathS)
					readBackD, errD := readbackClient.dc.FolderRead(context.Background(), pathD)
					assert.NoError(t, errS)
					assert.NoError(t, errD)
					TrimPrefixMap(readBackS, pathS)
					TrimPrefixMap(readBackD, pathD)

					if tt.wantNilSrc {
						assert.Nil(t, readBackS)
					} else {
						assert.Equal(t, orig, readBackS)
					}
					if tt.wantNilDst {
						assert.Nil(t, readBackD)
					} else {
						assert.Equal(t, orig, readBackD)
					}
				}
			}
		})
	}
}

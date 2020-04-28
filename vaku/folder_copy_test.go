package vaku

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderCopy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		giveSrc        string
		giveDst        string
		giveSrcLogical logical
		giveDstLogical logical
		giveOptions    []Option
		wantErr        []error
		wantNilDst     bool
	}{
		{
			name:    "copy one",
			giveSrc: "test/foo",
			giveDst: "copyone/test",
			wantErr: nil,
		},
		{
			name:    "copy all",
			giveSrc: "test",
			giveDst: "copyall/test",
			wantErr: nil,
		},
		{
			name:        "copy all absolute path",
			giveSrc:     "test",
			giveDst:     "copyallabs/test",
			giveOptions: []Option{WithabsolutePath(true)},
			wantErr:     nil,
		},
		{
			name:    "read err",
			giveSrc: "test/inner/again",
			giveDst: "readerr/test",
			giveSrcLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:    []error{ErrFolderCopy, ErrFolderRead, ErrFolderReadChan, ErrPathRead, ErrVaultRead},
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
			wantErr:    []error{ErrFolderCopy, ErrFolderWrite, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
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
					rbClient := cloneCLient(t, c)
					updateLogical(t, c, tt.giveSrcLogical, tt.giveDstLogical)

					pathS := addMountToPath(t, tt.giveSrc, ver[0])
					pathD := addMountToPath(t, tt.giveDst, ver[1])

					err := c.FolderCopy(context.Background(), pathS, pathD)
					compareErrors(t, err, tt.wantErr)

					readBackS, errS := rbClient.FolderRead(context.Background(), pathS)
					readBackD, errD := rbClient.dc.FolderRead(context.Background(), pathD)
					assert.NoError(t, errS)
					assert.NoError(t, errD)

					if tt.wantNilDst {
						assert.Nil(t, readBackD)
					} else {
						TrimPrefixMap(readBackS, pathS)
						TrimPrefixMap(readBackD, pathD)
						assert.Equal(t, readBackS, readBackD)
					}
				}
			}
		})
	}
}

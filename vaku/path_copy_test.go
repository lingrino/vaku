package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathCopy(t *testing.T) {
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
			name:    "copy",
			giveSrc: "test/foo",
			giveDst: "copy/test/foo",
			wantErr: nil,
		},
		{
			name:    "overwrite",
			giveSrc: "test/foo",
			giveDst: "test/value",
			wantErr: nil,
		},
		{
			name:       "bad src mount",
			giveSrc:    noMountPrefix,
			giveDst:    "copy/test/foo",
			wantErr:    []error{ErrPathCopy, ErrPathWrite, ErrNilData},
			wantNilDst: true,
		},
		{
			name:       "bad dst mount",
			giveSrc:    "test/foo",
			giveDst:    noMountPrefix,
			wantErr:    []error{ErrPathCopy, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
		{
			name:    "inject read",
			giveSrc: "test/foo",
			giveDst: "copy/injectread",
			giveSrcLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:    []error{ErrPathCopy, ErrPathRead, ErrVaultRead},
			wantNilDst: true,
		},
		{
			name:    "inject write",
			giveSrc: "test/foo",
			giveDst: "copy/injectwrite",
			giveDstLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantErr:    []error{ErrPathCopy, ErrPathWrite, ErrVaultWrite},
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
					readbackClient := cloneCLient(t, c)
					updateLogical(t, c, tt.giveSrcLogical, tt.giveDstLogical)

					pathS := addMountToPath(t, tt.giveSrc, ver[0])
					pathD := addMountToPath(t, tt.giveDst, ver[1])

					err := c.PathCopy(pathS, pathD)
					compareErrors(t, err, tt.wantErr)

					readBackS, errS := readbackClient.PathRead(pathS)
					readBackD, errD := readbackClient.pathReadDst(pathD)
					assert.NoError(t, errS)
					assert.NoError(t, errD)

					if tt.wantNilDst {
						assert.Nil(t, readBackD)
					} else {
						assert.Equal(t, readBackS, readBackD)
					}
				}
			}
		})
	}
}

package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMove(t *testing.T) {
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
			name:       "move",
			giveSrc:    "test/foo",
			giveDst:    "move/test/foo",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			name:       "overwrite",
			giveSrc:    "test/foo",
			giveDst:    "test/value",
			wantErr:    nil,
			wantNilSrc: true,
		},
		{
			name:       "bad src mount",
			giveSrc:    noMountPrefix,
			giveDst:    "move/test/foo",
			wantErr:    []error{ErrPathMove, ErrPathCopy, ErrPathWrite, ErrNilData},
			wantNilDst: true,
		},
		{
			name:       "bad dst mount",
			giveSrc:    "test/foo",
			giveDst:    noMountPrefix,
			wantErr:    []error{ErrPathMove, ErrPathCopy, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
		{
			name:    "inject read",
			giveSrc: "test/foo",
			giveDst: "move/injectread",
			giveSrcLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:    []error{ErrPathMove, ErrPathCopy, ErrPathRead, ErrVaultRead},
			wantNilDst: true,
		},
		{
			name:    "inject write",
			giveSrc: "test/foo",
			giveDst: "move/injectwrite",
			giveDstLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantErr:    []error{ErrPathMove, ErrPathCopy, ErrPathWrite, ErrVaultWrite},
			wantNilDst: true,
		},
		{
			name:    "inject delete",
			giveSrc: "test/foo",
			giveDst: "move/injectdelete",
			giveSrcLogical: &errLogical{
				err: errInject,
				op:  "Delete",
			},
			wantErr:    []error{ErrPathMove, ErrPathDelete, ErrVaultDelete},
			wantNilSrc: false,
			wantNilDst: false,
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

					orig, err := rbClient.PathRead(pathS)
					assert.NoError(t, err)

					err = c.PathMove(pathS, pathD)
					compareErrors(t, err, tt.wantErr)

					readBackS, errS := rbClient.PathRead(pathS)
					readBackD, errD := rbClient.dc.PathRead(pathD)
					assert.NoError(t, errS)
					assert.NoError(t, errD)

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

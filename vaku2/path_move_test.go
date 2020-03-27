package vaku2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		giveSource        string
		giveDest          string
		giveSourceLogical logical
		giveDestLogical   logical
		giveOptions       []Option
		wantErr           []error
		wantNilSource     bool
		wantNilDest       bool
	}{
		{
			name:          "move",
			giveSource:    "test/foo",
			giveDest:      "move/test/foo",
			wantErr:       nil,
			wantNilSource: true,
		},
		{
			name:          "overwrite",
			giveSource:    "test/foo",
			giveDest:      "test/value",
			wantErr:       nil,
			wantNilSource: true,
		},
		{
			name:        "bad source mount",
			giveSource:  noMountPrefix,
			giveDest:    "move/test/foo",
			wantErr:     []error{ErrPathMove, ErrPathCopy, ErrVaultWrite},
			wantNilDest: true,
		},
		{
			name:        "bad dest mount",
			giveSource:  "test/foo",
			giveDest:    noMountPrefix,
			wantErr:     []error{ErrPathMove, ErrPathCopy, ErrVaultWrite},
			wantNilDest: true,
		},
		{
			name:       "inject read",
			giveSource: "test/foo",
			giveDest:   "move/injectread",
			giveSourceLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:     []error{ErrPathMove, ErrPathCopy, ErrVaultRead},
			wantNilDest: true,
		},
		{
			name:       "inject write",
			giveSource: "test/foo",
			giveDest:   "move/injectwrite",
			giveDestLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantErr:     []error{ErrPathMove, ErrPathCopy, ErrVaultWrite},
			wantNilDest: true,
		},
		{
			name:       "inject delete",
			giveSource: "test/foo",
			giveDest:   "move/injectdelete",
			giveSourceLogical: &errLogical{
				err: errInject,
				op:  "Delete",
			},
			wantErr:       []error{ErrPathMove, ErrVaultDelete},
			wantNilSource: false,
			wantNilDest:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			for _, ver := range versionProduct {
				ln, client := testClient(t, tt.giveOptions...)
				defer ln.Close()

				lnS, lnD, clientDD := testClientDiffDest(t, tt.giveOptions...)
				defer lnS.Close()
				defer lnD.Close()

				for _, c := range []*Client{client, clientDD} {
					readbackClient := cloneCLient(t, c)
					updateLogical(t, c, tt.giveSourceLogical, tt.giveDestLogical)

					pathS := addMountToPath(t, tt.giveSource, ver[0])
					pathD := addMountToPath(t, tt.giveDest, ver[1])

					orig, err := readbackClient.PathRead(pathS)
					assert.NoError(t, err)

					err = c.PathMove(pathS, pathD)
					compareErrors(t, err, tt.wantErr)

					readBackS, errS := readbackClient.PathRead(pathS)
					readBackD, errD := readbackClient.PathReadDest(pathD)
					assert.NoError(t, errS)
					assert.NoError(t, errD)

					if tt.wantNilSource {
						assert.Nil(t, readBackS)
					} else {
						assert.Equal(t, orig, readBackS)
					}
					if tt.wantNilDest {
						assert.Nil(t, readBackD)
					} else {
						assert.Equal(t, orig, readBackD)
					}
				}
			}
		})
	}
}

package vaku2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathCopy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		giveSource        string
		giveDest          string
		giveSourceLogical logical
		giveDestLogical   logical
		giveOptions       []Option
		wantErr           []error
		wantNilDest       bool
	}{
		{
			name:       "copy",
			giveSource: "test/foo",
			giveDest:   "copy/test/foo",
			wantErr:    nil,
		},
		{
			name:       "overwrite",
			giveSource: "test/foo",
			giveDest:   "test/value",
			wantErr:    nil,
		},
		{
			name:        "bad source mount",
			giveSource:  noMountPrefix,
			giveDest:    "copy/test/foo",
			wantErr:     []error{ErrPathCopy, ErrVaultWrite},
			wantNilDest: true,
		},
		{
			name:        "bad dest mount",
			giveSource:  "test/foo",
			giveDest:    noMountPrefix,
			wantErr:     []error{ErrPathCopy, ErrVaultWrite},
			wantNilDest: true,
		},
		{
			name:       "inject read",
			giveSource: "test/foo",
			giveDest:   "copy/injectread",
			giveSourceLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:     []error{ErrPathCopy, ErrVaultRead},
			wantNilDest: true,
		},
		{
			name:       "inject write",
			giveSource: "test/foo",
			giveDest:   "copy/injectwrite",
			giveDestLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantErr:     []error{ErrPathCopy, ErrVaultWrite},
			wantNilDest: true,
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

					err := c.PathCopy(pathS, pathD)
					compareErrors(t, err, tt.wantErr)

					readBackS, errS := readbackClient.PathRead(pathS)
					readBackD, errD := readbackClient.PathReadDest(pathD)
					assert.NoError(t, errS)
					assert.NoError(t, errD)

					if tt.wantNilDest {
						assert.Nil(t, readBackD)
					} else {
						assert.Equal(t, readBackS, readBackD)
					}
				}
			}
		})
	}
}

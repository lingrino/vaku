package vaku2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathCopy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveSource  string
		giveDest    string
		giveLogical logical
		giveOptions []Option
		wantErr     error
		wantNilDest bool
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
			wantErr:     ErrPathCopy,
			wantNilDest: true,
		},
		{
			name:        "bad dest mount",
			giveSource:  "test/foo",
			giveDest:    noMountPrefix,
			wantErr:     ErrPathCopy,
			wantNilDest: true,
		},
		{
			name:       "inject read",
			giveSource: "test/foo",
			giveDest:   "copy/injectread",
			giveLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:     ErrPathCopy,
			wantNilDest: true,
		},
		{
			name:       "inject write",
			giveSource: "test/foo",
			giveDest:   "copy/injectwrite",
			giveLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantErr:     ErrPathCopy,
			wantNilDest: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			versionProduct := [][2]string{
				{"1", "1"},
				{"2", "2"},
				{"1", "2"},
				{"2", "1"},
			}

			for _, ver := range versionProduct {
				ln, client := testClient(t, tt.giveOptions...)
				defer ln.Close()
				readbackClient := cloneCLient(t, client)
				updateLogical(t, client, tt.giveLogical)

				pathS := addMountToPath(t, tt.giveSource, ver[0])
				pathD := addMountToPath(t, tt.giveDest, ver[1])

				err := client.PathCopy(pathS, pathD)
				assert.True(t, errors.Is(err, tt.wantErr), err)

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
		})
	}
}

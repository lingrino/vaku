package vaku2

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathMove(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		giveSource    string
		giveDest      string
		giveLogical   logical
		giveOptions   []Option
		wantErr       error
		wantNilSource bool
		wantNilDest   bool
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
			wantErr:     ErrPathMove,
			wantNilDest: true,
		},
		{
			name:        "bad dest mount",
			giveSource:  "test/foo",
			giveDest:    noMountPrefix,
			wantErr:     ErrPathMove,
			wantNilDest: true,
		},
		{
			name:       "inject read",
			giveSource: "test/foo",
			giveDest:   "move/injectread",
			giveLogical: &errLogical{
				err: errInject,
				op:  "Read",
			},
			wantErr:     ErrPathMove,
			wantNilDest: true,
		},
		{
			name:       "inject write",
			giveSource: "test/foo",
			giveDest:   "move/injectwrite",
			giveLogical: &errLogical{
				err: errInject,
				op:  "Write",
			},
			wantErr:     ErrPathMove,
			wantNilDest: true,
		},
		{
			name:       "inject delete",
			giveSource: "test/foo",
			giveDest:   "move/injectdelete",
			giveLogical: &errLogical{
				err: errInject,
				op:  "Delete",
			},
			wantErr:       ErrPathMove,
			wantNilSource: false,
			wantNilDest:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Tests with the same source/destination
			ln, client := testClient(t, tt.giveOptions...)
			defer ln.Close()
			readbackClient := cloneCLient(t, client)

			updateLogical(t, client, tt.giveLogical)

			for _, ver := range kvMountVersions {
				pathS := addMountToPath(t, tt.giveSource, ver)
				pathD := addMountToPath(t, tt.giveDest, ver)

				origS, err := readbackClient.PathRead(pathS)
				assert.NoError(t, err)

				err = client.PathMove(pathS, pathD)
				assert.True(t, errors.Is(err, tt.wantErr), err)

				readBackS, errS := readbackClient.PathRead(pathS)
				readBackD, errD := readbackClient.PathRead(pathD)
				assert.NoError(t, errS)
				assert.NoError(t, errD)

				if tt.wantNilSource {
					assert.Nil(t, readBackS)
				} else {
					assert.Equal(t, origS, readBackS)
				}
				if tt.wantNilDest {
					assert.Nil(t, readBackD)
				} else {
					assert.Equal(t, origS, readBackD)
				}
			}

			// Tests with different source/destination
			lnS, lnD, clientD := testClientDiffDest(t, tt.giveOptions...)
			defer lnS.Close()
			defer lnD.Close()

			updateLogical(t, clientD, tt.giveLogical)

			for _, ver := range kvMountVersions {
				pathS := addMountToPath(t, tt.giveSource, ver)
				pathD := addMountToPath(t, tt.giveDest, ver)

				origS, err := readbackClient.PathRead(pathS)
				assert.NoError(t, err)

				err = clientD.PathMove(pathS, pathD)
				assert.True(t, errors.Is(err, tt.wantErr), err)

				readBackS, errS := readbackClient.PathRead(pathS)
				readBackD, errD := readbackClient.PathRead(pathD)
				assert.NoError(t, errS)
				assert.NoError(t, errD)

				if tt.wantNilSource {
					assert.Nil(t, readBackS)
				} else {
					assert.Equal(t, origS, readBackS)
				}
				if tt.wantNilDest {
					assert.Nil(t, readBackD)
				} else {
					assert.Equal(t, origS, readBackD)
				}
			}

			// 	// Tests with different source/destination and ver 1 -> 2
			// 	lnS, lnD, client = testClientDiffDest(t, tt.giveOptions...)
			// 	defer lnS.Close()
			// 	defer lnD.Close()

			// 	updateLogical(t, client, tt.giveLogical)

			// 	pathS := addMountToPath(t, tt.giveSource, "1")
			// 	pathD := addMountToPath(t, tt.giveDest, "2")

			// 	origS, err := readbackClient.PathRead(pathS)
			// 	assert.NoError(t, err)

			// 	err = client.PathMove(pathS, pathD)
			// 	assert.True(t, errors.Is(err, tt.wantErr), err)

			// 	readBackS, errS := readbackClient.PathRead(pathS)
			// 	readBackD, errD := readbackClient.PathRead(pathD)
			// 	assert.NoError(t, errS)
			// 	assert.NoError(t, errD)

			// 	if tt.wantNilSource {
			// 		assert.Nil(t, readBackS)
			// 	} else {
			// 		assert.Equal(t, origS, readBackS)
			// 	}
			// 	if tt.wantNilDest {
			// 		assert.Nil(t, readBackD)
			// 	} else {
			// 		assert.Equal(t, origS, readBackD)
			// 	}

			// 	// Tests with different source/destination and ver 2 -> 1
			// 	lnS, lnD, client = testClientDiffDest(t, tt.giveOptions...)
			// 	defer lnS.Close()
			// 	defer lnD.Close()

			// 	updateLogical(t, client, tt.giveLogical)

			// 	pathS = addMountToPath(t, tt.giveSource, "2")
			// 	pathD = addMountToPath(t, tt.giveDest, "1")

			// 	origS, err = readbackClient.PathRead(pathS)
			// 	assert.NoError(t, err)

			// 	err = client.PathMove(pathS, pathD)
			// 	assert.True(t, errors.Is(err, tt.wantErr), err)

			// 	readBackS, errS = readbackClient.PathRead(pathS)
			// 	readBackD, errD = readbackClient.PathRead(pathD)
			// 	assert.NoError(t, errS)
			// 	assert.NoError(t, errD)

			// 	if tt.wantNilSource {
			// 		assert.Nil(t, readBackS)
			// 	} else {
			// 		assert.Equal(t, origS, readBackS)
			// 	}
			// 	if tt.wantNilDest {
			// 		assert.Nil(t, readBackD)
			// 	} else {
			// 		assert.Equal(t, origS, readBackD)
			// 	}
		})
	}
}

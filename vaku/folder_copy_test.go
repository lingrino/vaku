package vaku

// import (
// 	"context"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// )

// func TestFolderCopy(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name           string
// 		giveSrc        string
// 		giveDst        string
// 		giveSrcLogical logical
// 		giveDstLogical logical
// 		giveOptions    []Option
// 		wantErr        []error
// 		wantNilDst     bool
// 	}{
// 		{
// 			name:    "copy one",
// 			giveSrc: "test/foo",
// 			giveDst: "copy/test/foo",
// 			wantErr: nil,
// 		},
// 		{
// 			name:    "copy all",
// 			giveSrc: "test",
// 			giveDst: "copy/test",
// 			wantErr: nil,
// 		},
// 	}

// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			for _, ver := range versionProduct {
// 				client := testClient(t, tt.giveOptions...)
// 				clientDD := testClientDiffDst(t, tt.giveOptions...)

// 				for _, c := range []*Client{client, clientDD} {
// 					readbackClient := cloneCLient(t, c)
// 					updateLogical(t, c, tt.giveSrcLogical, tt.giveDstLogical)

// 					pathS := addMountToPath(t, tt.giveSrc, ver[0])
// 					pathD := addMountToPath(t, tt.giveDst, ver[1])

// 					err := c.FolderCopy(context.Background(), pathS, pathD)
// 					compareErrors(t, err, tt.wantErr)

// 					readBackS, errS := readbackClient.FolderRead(context.Background(), pathS)
// 					readBackD, errD := readbackClient.folderReadDst(context.Background(), pathD)
// 					assert.NoError(t, errS)
// 					assert.NoError(t, errD)

// 					if tt.wantNilDst {
// 						assert.Nil(t, readBackD)
// 					} else {
// 						assert.Equal(t, readBackS, readBackD)
// 					}
// 				}
// 			}
// 		})
// 	}
// }

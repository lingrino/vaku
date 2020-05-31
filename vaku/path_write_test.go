package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathWrite(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		give           string
		giveData       map[string]interface{}
		wantErr        []error
		wantNoReadback bool
	}{
		{
			give: "newpath",
			giveData: map[string]interface{}{
				"0": "1",
				"2": "3",
				"4": "5",
			},
			wantErr: nil,
		},
		{
			give: "0/1",
			giveData: map[string]interface{}{
				"100": "200",
			},
			wantErr: nil,
		},
		{
			give:     "nildata",
			giveData: nil,
			wantErr:  []error{ErrPathWrite, ErrNilData},
		},
		{
			give:           "injecterror",
			giveData:       map[string]interface{}{"0": "1"},
			wantErr:        []error{ErrPathWrite, ErrVaultWrite},
			wantNoReadback: true,
		},
		{
			give:           mountless,
			giveData:       map[string]interface{}{"0": "1"},
			wantErr:        []error{ErrPathWrite, ErrRewritePath, ErrMountInfo, ErrNoMount},
			wantNoReadback: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.give, func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPath(t, tt.give) {
				prefix := prefix
				t.Run(prefix, func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.PathWrite(PathJoin(prefix, tt.give), tt.giveData)
					compareErrors(t, err, tt.wantErr)

					if !tt.wantNoReadback {
						readBack, err := sharedReadBack.PathRead(PathJoin(prefix, tt.give))
						assert.NoError(t, err)
						assert.Equal(t, tt.giveData, readBack)
					}
				})
			}
		})
	}
}

package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathUpate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give           string
		giveData       map[string]interface{}
		wantData       map[string]interface{}
		wantErr        []error
		wantNoReadback bool
	}{
		{
			give:     "newpath",
			giveData: map[string]interface{}{"0": "1"},
			wantData: map[string]interface{}{"0": "1"},
			wantErr:  nil,
		},
		{
			give: "0/1",
			giveData: map[string]interface{}{
				"100": "101",
			},
			wantData: map[string]interface{}{
				"2":   "3",
				"100": "101",
			},
			wantErr: nil,
		},
		{
			give:     "nildata",
			giveData: nil,
			wantErr:  []error{ErrPathUpdate, ErrNilData},
		},
		{
			give:     "0/4/5",
			giveData: nil,
			wantData: map[string]interface{}{
				"6": "7",
			},
			wantErr: []error{ErrPathUpdate, ErrNilData},
		},
		{
			give: mountless,
			giveData: map[string]interface{}{
				"0": "1",
			},
			wantErr:        []error{ErrPathUpdate, ErrPathRead, ErrRewritePath, ErrMountInfo, ErrNoMount},
			wantNoReadback: true,
		},
		{
			give: "error/write/inject",
			giveData: map[string]interface{}{
				"0": "1",
			},
			wantErr:        []error{ErrPathUpdate, ErrPathWrite, ErrVaultWrite},
			wantNoReadback: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				prefix := prefix
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.PathUpdate(PathJoin(prefix, tt.give), tt.giveData)
					compareErrors(t, err, tt.wantErr)

					if !tt.wantNoReadback {
						readBack, err := sharedVakuClean.PathRead(PathJoin(prefix, tt.give))
						assert.NoError(t, err)
						assert.Equal(t, tt.wantData, readBack)
					}
				})
			}
		})
	}
}

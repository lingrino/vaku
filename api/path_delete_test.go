package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give           string
		wantErr        []error
		wantNoReadback bool
	}{
		{
			give:    "0/1",
			wantErr: nil,
		},
		{
			give:    "fake",
			wantErr: nil,
		},
		{
			give:           mountless,
			wantErr:        []error{ErrPathDelete, ErrRewritePath, ErrMountInfo, ErrNoMount},
			wantNoReadback: true,
		},
		{
			give:           "error/delete/inject",
			wantErr:        []error{ErrPathDelete, ErrVaultDelete},
			wantNoReadback: true,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					err := sharedVaku.PathDelete(PathJoin(prefix, tt.give))
					compareErrors(t, err, tt.wantErr)

					if !tt.wantNoReadback {
						readBack, err := sharedVakuClean.PathRead(PathJoin(prefix, tt.give))
						assert.NoError(t, err)
						assert.Nil(t, readBack)
					}
				})
			}
		})
	}
}

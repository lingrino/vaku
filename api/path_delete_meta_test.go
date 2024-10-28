package vaku

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathDeleteMeta(t *testing.T) {
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
			give:           "error/delete/inject",
			wantErr:        []error{ErrPathDeleteMeta, ErrVaultDelete},
			wantNoReadback: true,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				if strings.HasPrefix(prefix, "kv1") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						err := sharedVaku.PathDeleteMeta(PathJoin(prefix, tt.give))
						compareErrors(t, err, []error{ErrPathDeleteMeta, ErrMountVersion})
					})
				}
				if strings.HasPrefix(prefix, "kv2") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						err := sharedVaku.PathDeleteMeta(PathJoin(prefix, tt.give))
						compareErrors(t, err, tt.wantErr)

						if !tt.wantNoReadback {
							readBack, err := sharedVakuClean.PathRead(PathJoin(prefix, tt.give))
							assert.NoError(t, err)
							assert.Nil(t, readBack)
						}
					})
				}
			}
		})
	}
}

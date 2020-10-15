package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		want    []string
		wantErr []error
	}{
		{
			give:    "0",
			want:    []string{"1", "4/"},
			wantErr: nil,
		},
		{
			give:    "0/4/13/24",
			want:    []string{"25/"},
			wantErr: nil,
		},
		{
			give:    "emptypath",
			want:    nil,
			wantErr: nil,
		},
		{
			give:    mountless,
			want:    nil,
			wantErr: []error{ErrPathList, ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			give:    "error/list/inject",
			want:    nil,
			wantErr: []error{ErrPathList, ErrVaultList},
		},
		{
			give:    "nildata/list/inject",
			want:    nil,
			wantErr: nil,
		},
		{
			give:    "nilkeys/list/inject",
			want:    nil,
			wantErr: []error{ErrPathList, ErrDecodeSecret},
		},
		{
			give:    "intkeys/list/inject",
			want:    nil,
			wantErr: []error{ErrPathList, ErrDecodeSecret},
		},
		{
			give:    "listintkeys/list/inject",
			want:    nil,
			wantErr: []error{ErrPathList, ErrDecodeSecret},
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

					list, err := sharedVaku.PathList(PathJoin(prefix, tt.give))
					TrimPrefixList(list, prefix)

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.want, list)
				})
			}
		})
	}
}

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
			give:    "fake",
			want:    nil,
			wantErr: nil,
		},
		{
			give:    mountless,
			want:    nil,
			wantErr: []error{ErrPathList, ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			give:    "injecterror",
			want:    nil,
			wantErr: []error{ErrPathList, ErrVaultList},
		},
		{
			give:    "injectdatanil",
			want:    nil,
			wantErr: nil,
		},
		{
			give:    "injectkeysnil",
			want:    nil,
			wantErr: []error{ErrPathList, ErrDecodeSecret},
		},
		{
			give:    "injectkeysint",
			want:    nil,
			wantErr: []error{ErrPathList, ErrDecodeSecret},
		},
		{
			give:    "injectkeyslistint",
			want:    nil,
			wantErr: []error{ErrPathList, ErrDecodeSecret},
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

					list, err := sharedVaku.PathList(PathJoin(prefix, tt.give))
					TrimPrefixList(list, prefix)

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.want, list)
				})
			}
		})
	}
}

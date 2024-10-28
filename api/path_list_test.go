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
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
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

func TestPathListIgnoreErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		want    []string
		wantErr []error
	}{
		{
			give:    "error/list/inject",
			want:    nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					client, err := NewClient(
						WithVaultSrcClient(testServer(t)),
						WithIgnoreAccessErrors(true),
					)
					assert.NoError(t, err)
					client.vl = &logicalInjector{realL: client.vl, t: t}

					read, err := client.PathList(PathJoin(prefix, tt.give))

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.want, read)
				})
			}
		})
	}
}

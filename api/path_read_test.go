package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		want    map[string]any
		wantErr []error
	}{
		{
			give: "0/1",
			want: map[string]any{
				"2": "3",
			},
			wantErr: nil,
		},
		{
			give: "0/4/13/17",
			want: map[string]any{
				"18": "19",
				"20": "21",
				"22": "23",
			},
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
			wantErr: []error{ErrPathRead, ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			give:    "error/read/inject",
			want:    nil,
			wantErr: []error{ErrPathRead, ErrVaultRead},
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

					read, err := sharedVaku.PathRead(PathJoin(prefix, tt.give))

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.want, read)
				})
			}
		})
	}
}

func TestExtractV2Read(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give map[string]any
		want map[string]any
	}{
		{
			give: nil,
			want: nil,
		},
		{
			give: map[string]any{"foo": "bar"},
			want: nil,
		},
		{
			give: map[string]any{"metadata": map[string]any{"foo": "bar"}},
			want: nil,
		},
		{
			give: map[string]any{"metadata": map[string]any{"deletion_time": ""}},
			want: nil,
		},
		{
			give: map[string]any{
				"metadata": map[string]any{
					"deletion_time": "",
					"destroyed":     false,
				},
			},
			want: nil,
		},
		{
			give: map[string]any{
				"metadata": map[string]any{
					"deletion_time": "",
					"destroyed":     false,
				},
				"data": map[string]any{
					"foo": "bar",
				},
			},
			want: map[string]any{
				"foo": "bar",
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := extractV2Read(tt.give)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestPathReadIgnoreErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		want    map[string]any
		wantErr []error
	}{
		{
			give:    "error/read/inject",
			want:    nil,
			wantErr: nil,
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

					client, err := NewClient(
						WithVaultSrcClient(testServer(t)),
						WithIgnoreAccessErrors(true),
					)
					assert.NoError(t, err)
					client.vl = &logicalInjector{realL: client.vl, t: t}

					read, err := client.PathRead(PathJoin(prefix, tt.give))

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.want, read)
				})
			}
		})
	}
}

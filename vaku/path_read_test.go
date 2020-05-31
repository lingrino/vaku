package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		want    map[string]interface{}
		wantErr []error
	}{
		{
			give: "0/1",
			want: map[string]interface{}{
				"2": "3",
			},
			wantErr: nil,
		},
		{
			give: "0/4/13/17",
			want: map[string]interface{}{
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
			give:    "injecterror",
			want:    nil,
			wantErr: []error{ErrPathRead, ErrVaultRead},
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
		give map[string]interface{}
		want map[string]interface{}
	}{
		{
			give: nil,
			want: nil,
		},
		{
			give: map[string]interface{}{"foo": "bar"},
			want: nil,
		},
		{
			give: map[string]interface{}{"metadata": map[string]interface{}{"foo": "bar"}},
			want: nil,
		},
		{
			give: map[string]interface{}{"metadata": map[string]interface{}{"deletion_time": ""}},
			want: nil,
		},
		{
			give: map[string]interface{}{
				"metadata": map[string]interface{}{
					"deletion_time": "",
					"destroyed":     false,
				},
			},
			want: nil,
		},
		{
			give: map[string]interface{}{
				"metadata": map[string]interface{}{
					"deletion_time": "",
					"destroyed":     false,
				},
				"data": map[string]interface{}{
					"foo": "bar",
				},
			},
			want: map[string]interface{}{
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

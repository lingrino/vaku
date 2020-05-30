package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathRead(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		give        string
		giveLogical logical
		giveOptions []Option
		want        map[string]interface{}
		wantErr     []error
	}{
		{
			name: "test/foo",
			give: "test/foo",
			want: map[string]interface{}{
				"value": "bar",
			},
			wantErr: nil,
		},
		{
			name: "test/inner/again/inner/UCrt6sZT",
			give: "test/inner/again/inner/UCrt6sZT",
			want: map[string]interface{}{
				"Eg5ljS7t": "6F1B5nBg",
				"quqr32S5": "81iY4HAN",
				"r6R0JUzX": "rs1mCRB5",
			},
			wantErr: nil,
		},
		{
			name:    "no secret",
			give:    "doesnotexist",
			want:    nil,
			wantErr: nil,
		},
		{
			name:    "no mount",
			give:    mountless,
			want:    nil,
			wantErr: []error{ErrPathRead, ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			name: "error",
			give: "test/foo",
			giveLogical: &errLogical{
				err: errInject,
			},
			want:    nil,
			wantErr: []error{ErrPathRead, ErrVaultRead},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client, _ := testSetup(t, tt.giveLogical, nil, tt.giveOptions...)

			for _, ver := range mountVersions {
				ver := ver
				t.Run(ver, func(t *testing.T) {
					t.Parallel()

					path := addMountToPath(t, tt.give, ver)

					read, err := client.PathRead(path)

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
			name: "nil",
			give: nil,
			want: nil,
		},
		{
			name: "no metadata",
			give: map[string]interface{}{"foo": "bar"},
			want: nil,
		},
		{
			name: "no deletion_time",
			give: map[string]interface{}{"metadata": map[string]interface{}{"foo": "bar"}},
			want: nil,
		},
		{
			name: "no destroyed",
			give: map[string]interface{}{"metadata": map[string]interface{}{"deletion_time": ""}},
			want: nil,
		},
		{
			name: "no data",
			give: map[string]interface{}{
				"metadata": map[string]interface{}{
					"deletion_time": "",
					"destroyed":     false,
				},
			},
			want: nil,
		},
		{
			name: "data",
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

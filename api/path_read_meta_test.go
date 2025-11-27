package vaku

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathReadMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		wantErr []error
		wantNil bool
	}{
		{
			give:    "0/1",
			wantErr: nil,
			wantNil: false,
		},
		{
			give:    "fake/path",
			wantErr: nil,
			wantNil: true,
		},
		{
			give:    "error/read/inject",
			wantErr: []error{ErrPathReadMeta, ErrVaultRead},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				if strings.HasPrefix(prefix, "kv1") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						_, err := sharedVaku.PathReadMeta(PathJoin(prefix, tt.give))
						compareErrors(t, err, []error{ErrPathReadMeta, ErrMountVersion})
					})
				}
				if strings.HasPrefix(prefix, "kv2") {
					t.Run(testName(prefix), func(t *testing.T) {
						t.Parallel()

						meta, err := sharedVaku.PathReadMeta(PathJoin(prefix, tt.give))
						compareErrors(t, err, tt.wantErr)

						if tt.wantNil {
							assert.Nil(t, meta)
						} else if tt.wantErr == nil {
							assert.NotNil(t, meta)
							assert.GreaterOrEqual(t, meta.CurrentVersion, 1)
							assert.NotEmpty(t, meta.Versions)
						}
					})
				}
			}
		})
	}
}

func TestExtractSecretMeta(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		give map[string]any
		want *SecretMeta
	}{
		{
			name: "full metadata with float64",
			give: map[string]any{
				"current_version": float64(3),
				"versions": map[string]any{
					"1": map[string]any{
						"created_time":  "2023-01-01T00:00:00Z",
						"deletion_time": "",
						"destroyed":     false,
					},
					"2": map[string]any{
						"created_time":  "2023-01-02T00:00:00Z",
						"deletion_time": "2023-01-03T00:00:00Z",
						"destroyed":     false,
					},
					"3": map[string]any{
						"created_time":  "2023-01-04T00:00:00Z",
						"deletion_time": "",
						"destroyed":     true,
					},
				},
			},
			want: &SecretMeta{
				CurrentVersion: 3,
				Versions: map[int]SecretVersionMeta{
					1: {CreatedTime: "2023-01-01T00:00:00Z", Deleted: false, Destroyed: false},
					2: {CreatedTime: "2023-01-02T00:00:00Z", Deleted: true, Destroyed: false},
					3: {CreatedTime: "2023-01-04T00:00:00Z", Deleted: false, Destroyed: true},
				},
			},
		},
		{
			name: "full metadata with json.Number",
			give: map[string]any{
				"current_version": json.Number("5"),
				"versions": map[string]any{
					"1": map[string]any{
						"created_time":  "2023-01-01T00:00:00Z",
						"deletion_time": "",
						"destroyed":     false,
					},
				},
			},
			want: &SecretMeta{
				CurrentVersion: 5,
				Versions: map[int]SecretVersionMeta{
					1: {CreatedTime: "2023-01-01T00:00:00Z", Deleted: false, Destroyed: false},
				},
			},
		},
		{
			name: "nil data",
			give: nil,
			want: &SecretMeta{
				CurrentVersion: 0,
				Versions:       map[int]SecretVersionMeta{},
			},
		},
		{
			name: "empty versions",
			give: map[string]any{
				"current_version": float64(0),
				"versions":        map[string]any{},
			},
			want: &SecretMeta{
				CurrentVersion: 0,
				Versions:       map[int]SecretVersionMeta{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := extractSecretMeta(tt.give)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

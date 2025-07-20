package vaku

import (
	"strings"
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
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
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

					read, err := client.PathRead(PathJoin(prefix, tt.give))

					compareErrors(t, err, tt.wantErr)
					assert.Equal(t, tt.want, read)
				})
			}
		})
	}
}

func TestPathReadMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		wantErr []error
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
			give:    mountless,
			wantErr: []error{ErrPathReadMetadata, ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			give:    "error/read/inject",
			wantErr: []error{ErrPathReadMetadata, ErrVaultRead},
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					// Only test on v2 mounts since metadata is v2 only
					if !strings.Contains(prefix, "v2") {
						t.Skip("PathReadMetadata only works on KV v2 mounts")
					}

					metadata, err := sharedVaku.PathReadMetadata(PathJoin(prefix, tt.give))

					compareErrors(t, err, tt.wantErr)
					if err == nil && metadata != nil {
						// Verify metadata structure for valid paths
						assert.Contains(t, metadata, "versions")
						assert.Contains(t, metadata, "current_version")
					}
				})
			}
		})
	}
}

func TestPathReadVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give    string
		version int
		wantErr []error
	}{
		{
			give:    "0/1",
			version: 1,
			wantErr: nil,
		},
		{
			give:    "fake",
			version: 1,
			wantErr: nil,
		},
		{
			give:    mountless,
			version: 1,
			wantErr: []error{ErrPathReadVersion, ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			give:    "fake/nonexistent",
			version: 1,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(testName(tt.give), func(t *testing.T) {
			t.Parallel()
			for _, prefix := range seededPrefixes(t, tt.give) {
				t.Run(testName(prefix), func(t *testing.T) {
					t.Parallel()

					// Only test on v2 mounts since versions are v2 only
					if !strings.Contains(prefix, "v2") {
						t.Skip("PathReadVersion only works on KV v2 mounts")
					}

					result, err := sharedVaku.PathReadVersion(PathJoin(prefix, tt.give), tt.version)

					compareErrors(t, err, tt.wantErr)
					if err == nil && result != nil {
						// For valid existing paths, we should get the same data as PathRead
						expected, errExpected := sharedVaku.PathRead(PathJoin(prefix, tt.give))
						assert.NoError(t, errExpected)
						assert.Equal(t, expected, result)
					}
				})
			}
		})
	}
}

func TestPathReadMetadataKV1(t *testing.T) {
	t.Parallel()

	// Test that PathReadMetadata returns proper error on KV v1
	for _, prefix := range seededPrefixes(t, "0/1") {
		t.Run(testName(prefix), func(t *testing.T) {
			t.Parallel()

			// Only test on v1 mounts
			if !strings.Contains(prefix, "v1") {
				t.Skip("Testing KV v1 error handling")
			}

			_, err := sharedVaku.PathReadMetadata(PathJoin(prefix, "0/1"))
			expectedErrors := []error{ErrPathReadMetadata, ErrMountVersion}
			compareErrors(t, err, expectedErrors)
		})
	}
}

func TestPathReadVersionKV1(t *testing.T) {
	t.Parallel()

	// Test that PathReadVersion returns proper error on KV v1
	for _, prefix := range seededPrefixes(t, "0/1") {
		t.Run(testName(prefix), func(t *testing.T) {
			t.Parallel()

			// Only test on v1 mounts
			if !strings.Contains(prefix, "v1") {
				t.Skip("Testing KV v1 error handling")
			}

			_, err := sharedVaku.PathReadVersion(PathJoin(prefix, "0/1"), 1)
			expectedErrors := []error{ErrPathReadVersion, ErrMountVersion}
			compareErrors(t, err, expectedErrors)
		})
	}
}

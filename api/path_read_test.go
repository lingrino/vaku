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

func TestPathReadVersion(t *testing.T) {
	t.Parallel()

	for _, prefix := range seededPrefixes(t, "readversion") {
		if strings.HasPrefix(prefix, "kv1") {
			t.Run(testName(prefix, "kv1 error"), func(t *testing.T) {
				t.Parallel()

				_, err := sharedVaku.PathReadVersion(PathJoin(prefix, "0/1"), 1)
				compareErrors(t, err, []error{ErrPathReadVersion, ErrMountVersion})
			})
		}

		if strings.HasPrefix(prefix, "kv2") {
			t.Run(testName(prefix, "basic read"), func(t *testing.T) {
				t.Parallel()

				srcPath := PathJoin(prefix, "readversion/basic")

				// Write multiple versions
				v1 := map[string]any{"version": "1"}
				v2 := map[string]any{"version": "2"}

				err := sharedVakuClean.PathWrite(srcPath, v1)
				assert.NoError(t, err)
				err = sharedVakuClean.PathWrite(srcPath, v2)
				assert.NoError(t, err)

				// Read version 1
				read1, err := sharedVaku.PathReadVersion(srcPath, 1)
				assert.NoError(t, err)
				assert.Equal(t, v1, read1)

				// Read version 2
				read2, err := sharedVaku.PathReadVersion(srcPath, 2)
				assert.NoError(t, err)
				assert.Equal(t, v2, read2)
			})

			t.Run(testName(prefix, "nonexistent version"), func(t *testing.T) {
				t.Parallel()

				srcPath := PathJoin(prefix, "readversion/nonexistent")

				// Write one version
				err := sharedVakuClean.PathWrite(srcPath, map[string]any{"foo": "bar"})
				assert.NoError(t, err)

				// Try to read version that doesn't exist
				read, err := sharedVaku.PathReadVersion(srcPath, 999)
				assert.NoError(t, err)
				assert.Nil(t, read)
			})

			t.Run(testName(prefix, "read error"), func(t *testing.T) {
				t.Parallel()

				_, err := sharedVaku.PathReadVersion(PathJoin(prefix, "error/read/inject"), 1)
				compareErrors(t, err, []error{ErrPathReadVersion, ErrVaultRead})
			})
		}
	}
}

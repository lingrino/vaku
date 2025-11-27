package vaku

import (
	"testing"

	"github.com/stretchr/testify/assert"

	vault "github.com/hashicorp/vault/api"
)

func TestStaticMountProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		version     string
		wantPath    string
		wantVersion string
		wantType    string
	}{
		{
			name:        "kv v2 mount",
			path:        "secret/",
			version:     "2",
			wantPath:    "secret/",
			wantVersion: "2",
			wantType:    "kv",
		},
		{
			name:        "kv v1 mount",
			path:        "kv1/",
			version:     "1",
			wantPath:    "kv1/",
			wantVersion: "1",
			wantType:    "kv",
		},
		{
			name:        "nested path",
			path:        "my/secret/",
			version:     "2",
			wantPath:    "my/secret/",
			wantVersion: "2",
			wantType:    "kv",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			provider := NewStaticMountProvider(tt.path, tt.version)
			mounts, err := provider.ListMounts()

			assert.NoError(t, err)
			assert.Len(t, mounts, 1)
			assert.Equal(t, tt.wantPath, mounts[0].Path)
			assert.Equal(t, tt.wantVersion, mounts[0].Version)
			assert.Equal(t, tt.wantType, mounts[0].Type)
		})
	}
}

func TestMountInfoWithStaticProvider(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		mountPath     string
		mountVersion  string
		queryPath     string
		wantMountPath string
		wantVersion   mountVersion
		wantErr       []error
	}{
		{
			name:          "matching v2 path",
			mountPath:     "secret/",
			mountVersion:  "2",
			queryPath:     "secret/foo/bar",
			wantMountPath: "secret/",
			wantVersion:   mv2,
		},
		{
			name:          "matching v1 path",
			mountPath:     "kv1/",
			mountVersion:  "1",
			queryPath:     "kv1/foo/bar",
			wantMountPath: "kv1/",
			wantVersion:   mv1,
		},
		{
			name:          "non-matching path",
			mountPath:     "secret/",
			mountVersion:  "2",
			queryPath:     "other/foo/bar",
			wantMountPath: "",
			wantVersion:   mv0,
			wantErr:       []error{ErrMountInfo, ErrNoMount},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			vc, err := vault.NewClient(vault.DefaultConfig())
			assert.NoError(t, err)

			provider := NewStaticMountProvider(tt.mountPath, tt.mountVersion)
			client, err := NewClient(
				WithVaultClient(vc),
				WithMountProvider(provider),
			)
			assert.NoError(t, err)

			path, vers, err := client.mountInfo(tt.queryPath)

			assert.Equal(t, tt.wantMountPath, path)
			assert.Equal(t, tt.wantVersion, vers)
			compareErrors(t, err, tt.wantErr)
		})
	}
}

func TestMountInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give        string
		wantPath    string
		wantVersion mountVersion
		wantErr     []error
	}{
		{
			give:        "nomount",
			wantPath:    "",
			wantVersion: mv0,
			wantErr:     []error{ErrMountInfo, ErrNoMount},
		},
		{
			give:        "sys/",
			wantPath:    "sys/",
			wantVersion: mv0,
		},
		{
			give:        "kv1/",
			wantPath:    "kv1/",
			wantVersion: mv1,
		},
		{
			give:        "kv2/",
			wantPath:    "kv2/",
			wantVersion: mv2,
		},
	}

	t.Run("empty client", func(t *testing.T) {
		t.Parallel()

		vc, err := vault.NewClient(vault.DefaultConfig())
		assert.NoError(t, err)

		client, err := NewClient(WithVaultClient(vc))
		assert.NoError(t, err)

		path, vers, err := client.mountInfo("kv0")

		assert.Empty(t, path)
		assert.Equal(t, mv0, vers)
		compareErrors(t, err, []error{ErrMountInfo, ErrListMounts})
	})

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(WithVaultClient(testServer(t)))
			assert.NoError(t, err)

			path, vers, err := client.mountInfo(tt.give)

			assert.Equal(t, tt.wantPath, path)
			assert.Equal(t, tt.wantVersion, vers)
			compareErrors(t, err, tt.wantErr)
		})
	}
}

func TestMountStringToVersion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give string
		want mountVersion
	}{
		{
			give: "---",
			want: mv0,
		},
		{
			give: "0",
			want: mv0,
		},
		{
			give: "1",
			want: mv1,
		},
		{
			give: "2",
			want: mv2,
		},
		{
			give: "3",
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			t.Parallel()

			mv := mountStringToVersion(tt.give)
			assert.Equal(t, tt.want, mv)
		})
	}
}

func TestRewritePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		give        string
		giveOp      vaultOperation
		wantPath    string
		wantVersion mountVersion
		wantErr     []error
	}{
		{
			give:        "nomount",
			giveOp:      vaultRead,
			wantPath:    "",
			wantVersion: mv0,
			wantErr:     []error{ErrRewritePath, ErrMountInfo, ErrNoMount},
		},
		{
			give:        "kv1/a/b/c",
			giveOp:      vaultList,
			wantPath:    "kv1/a/b/c",
			wantVersion: mv1,
		},
		{
			give:        "kv2/a/b/c",
			giveOp:      vaultList,
			wantPath:    "kv2/metadata/a/b/c",
			wantVersion: mv2,
		},
		{
			give:        "kv1/a/b/c",
			giveOp:      vaultRead,
			wantPath:    "kv1/a/b/c",
			wantVersion: mv1,
		},
		{
			give:        "kv2/a/b/c",
			giveOp:      vaultRead,
			wantPath:    "kv2/data/a/b/c",
			wantVersion: mv2,
		},
		{
			give:        "kv1/a/b/c",
			giveOp:      vaultWrite,
			wantPath:    "kv1/a/b/c",
			wantVersion: mv1,
		},
		{
			give:        "kv2/a/b/c",
			giveOp:      vaultWrite,
			wantPath:    "kv2/data/a/b/c",
			wantVersion: mv2,
		},
		{
			give:        "kv1/a/b/c",
			giveOp:      vaultDelete,
			wantPath:    "kv1/a/b/c",
			wantVersion: mv1,
		},
		{
			give:        "kv2/a/b/c",
			giveOp:      vaultDelete,
			wantPath:    "kv2/data/a/b/c",
			wantVersion: mv2,
		},
		{
			give:        "kv1/a/b/c",
			giveOp:      vaultDestroy,
			wantPath:    "",
			wantVersion: mv1,
			wantErr:     []error{ErrMountVersion},
		},
		{
			give:        "kv2/a/b/c",
			giveOp:      vaultDestroy,
			wantPath:    "kv2/destroy/a/b/c",
			wantVersion: mv2,
		},
		{
			give:        "kv1/a/b/c",
			giveOp:      vaultDeleteMeta,
			wantPath:    "",
			wantVersion: mv1,
			wantErr:     []error{ErrMountVersion},
		},
		{
			give:        "kv2/a/b/c",
			giveOp:      vaultDeleteMeta,
			wantPath:    "kv2/metadata/a/b/c",
			wantVersion: mv2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			t.Parallel()

			client, err := NewClient(WithVaultClient(testServer(t)))
			assert.NoError(t, err)

			path, vers, err := client.rewritePath(tt.give, tt.giveOp)

			assert.Equal(t, tt.wantPath, path)
			assert.Equal(t, tt.wantVersion, vers)
			compareErrors(t, err, tt.wantErr)
		})
	}
}

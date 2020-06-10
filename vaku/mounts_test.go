package vaku

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

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

		vc, err := api.NewClient(api.DefaultConfig())
		assert.NoError(t, err)

		client, err := NewClient(WithVaultClient(vc))
		assert.NoError(t, err)

		path, vers, err := client.mountInfo("kv0")

		assert.Empty(t, path)
		assert.Equal(t, mv0, vers)
		compareErrors(t, err, []error{ErrMountInfo, ErrListMounts})
	})

	for _, tt := range tests {
		tt := tt
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
		tt := tt
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
			giveOp:      vaultDestroy,
			wantPath:    "kv1/a/b/c",
			wantVersion: mv1,
		},
		{
			give:        "kv2/a/b/c",
			giveOp:      vaultDestroy,
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
	}

	for _, tt := range tests {
		tt := tt
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

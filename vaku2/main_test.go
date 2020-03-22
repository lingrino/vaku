package vaku2

import (
	"errors"
	"net"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
)

var (
	// errInject is used when injecting errors in tests
	errInject = errors.New("injected error")
)

// When tests are looping over kvMountVersions and the path is noMountPrefix they will not prefix
// the path with the mount version to allow testing on a nonexistent mount.
var noMountPrefix = "nomount"

// kvMountVersions lists the types of kv mounts for vault. There are currently two k/v mount types
// and vaku supports both. Tests should run against each version and return the same results.
var kvMountVersions = []string{"1", "2"}

// seeds holds the canonical secret seeds for every test.
var seeds = map[string]map[string]interface{}{
	"test/foo": {
		"value": "bar",
	},
	"test/value": {
		"fizz": "buzz",
		"foo":  "bar",
	},
	"test/fizz": {
		"fizz": "buzz",
		"foo":  "bar",
	},
	"test/HToOeKKD": {
		"3zqxVbJY": "TvOjGxvC",
	},
	"test/inner/WKNC3muM": {
		"IY1C148K": "JxBfEt91",
		"iwVzPqbY": "0NH9GlR1",
	},
	"test/inner/A2xlzTfE": {
		"Eg5ljS7t": "BHRMKjj1",
		"quqr32S5": "pcidzSMW",
	},
	"test/inner/again/inner/UCrt6sZT": {
		"Eg5ljS7t": "6F1B5nBg",
		"quqr32S5": "81iY4HAN",
		"r6R0JUzX": "rs1mCRB5",
	},
}

// testServer creates a new inmem Vault server and returns a seeded client that points to it.
func testServer(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	core, _, token := vault.TestCoreUnsealedWithConfig(t, &vault.CoreConfig{Logger: hclog.NewNullLogger()})
	ln, addr := http.TestServer(t, core)

	client, err := api.NewClient(api.DefaultConfig())
	assert.NoError(t, err)

	client.SetToken(token)
	client.SetAddress(addr)

	for _, ver := range kvMountVersions {
		err = client.Sys().Mount(ver+"/", &api.MountInput{
			Type: "kv",
			Options: map[string]string{
				"version": ver,
			},
		})
		assert.NoError(t, err)

		for path, secret := range seeds {
			_, err := client.Logical().Write(PathJoin(ver, path), secret)
			assert.NoError(t, err)
		}
	}

	return ln, client
}

// errLogical implements logical and injects ouputs
type errLogical struct {
	secret *api.Secret
	err    error
}

func (e *errLogical) List(path string) (*api.Secret, error) {
	return e.secret, e.err
}

func (e *errLogical) Read(path string) (*api.Secret, error) {
	return e.secret, e.err
}

func (e *errLogical) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	return e.secret, e.err
}

func (e *errLogical) Delete(path string) (*api.Secret, error) {
	return e.secret, e.err
}

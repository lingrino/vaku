package vaku2

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
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

// versionProduct is all possible to/from version combinations for testing functions that should
// work across multiple mount versions.
var versionProduct = [4][2]string{
	{"1", "1"},
	{"2", "2"},
	{"1", "2"},
	{"2", "1"},
}

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

// It is necessary to to set up a TestMain here because vault.TestCoreUnsealedWithConfig() calls a
// function further down that uses a logger with a default config intead of the config passed in.
func TestMain(m *testing.M) {
	hclog.DefaultOutput = ioutil.Discard
	os.Exit(m.Run())
}

// testServer creates a new Vault server and returns a Vault API client that points to it.
func testServer(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	core, _, token := vault.TestCoreUnsealedWithConfig(t, &vault.CoreConfig{
		Logger: hclog.NewNullLogger(),
	})
	ln, addr := http.TestServer(t, core)

	apiClient, err := api.NewClient(api.DefaultConfig())
	assert.NoError(t, err)

	apiClient.SetToken(token)
	err = apiClient.SetAddress(addr)
	assert.NoError(t, err)

	return ln, apiClient
}

// testServerSeeded creates a seeded Vault server and returns an API client that points to it.
func testServerSeeded(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	ln, client := testServer(t)

	for _, ver := range kvMountVersions {
		err := client.Sys().Mount(ver+"/", &api.MountInput{
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

// testClient returns a client that points to an seeded server.
func testClient(t *testing.T, opts ...Option) (net.Listener, *Client) {
	t.Helper()

	ln, apiClient := testServerSeeded(t)

	client, err := NewClient(append(opts, WithVaultClient(apiClient))...)
	assert.NoError(t, err)

	return ln, client
}

// testClientDiffDest returns a client that points source and dest at different seeded servers.
func testClientDiffDest(t *testing.T, opts ...Option) (net.Listener, net.Listener, *Client) {
	t.Helper()

	ln, apiClientS := testServerSeeded(t)
	lnD, apiClientD := testServerSeeded(t)

	client, err := NewClient(append(opts,
		WithVaultSourceClient(apiClientS),
		WithVaultDestClient(apiClientD),
	)...)
	assert.NoError(t, err)

	return ln, lnD, client
}

// cloneCLient cpies a client. Don't use this outside of tests.
func cloneCLient(t *testing.T, c *Client) *Client {
	t.Helper()

	cpy := *c
	return &cpy
}

// errLogical implements logical and injects ouputs.
type errLogical struct {
	secret *api.Secret
	err    error

	// if op != "" all functions will pass to the real client except the one named in op
	op      string
	realL   logical
	useDest bool
}

func (e *errLogical) Delete(path string) (*api.Secret, error) {
	if e.op != "Delete" && e.op != "" {
		return e.realL.Delete(path)
	}
	return e.secret, e.err
}

func (e *errLogical) List(path string) (*api.Secret, error) {
	if e.op != "List" && e.op != "" {
		return e.realL.List(path)
	}
	return e.secret, e.err
}

func (e *errLogical) Read(path string) (*api.Secret, error) {
	if e.op != "Read" && e.op != "" {
		return e.realL.Read(path)
	}
	return e.secret, e.err
}

func (e *errLogical) Write(path string, data map[string]interface{}) (*api.Secret, error) {
	if e.op != "Write" && e.op != "" {
		return e.realL.Write(path, data)
	}
	return e.secret, e.err
}

// updateLogical is used in tests with tt.giveSourceLogical.
func updateLogical(t *testing.T, c *Client, sourceL logical, destL logical) {
	t.Helper()

	if sourceL != nil {
		sl, ok := sourceL.(*errLogical)
		if ok {
			sl.realL = c.sourceL
			c.sourceL = sl
		} else {
			c.sourceL = sourceL
		}
	}
	if destL != nil {
		dl, ok := destL.(*errLogical)
		if ok {
			dl.realL = c.destL
			c.destL = dl
		} else {
			c.destL = destL
		}
	}
}

// addMountToPath prefixes a path with a mount if path is not the special noMountPrefix.
func addMountToPath(t *testing.T, path string, mount string) string {
	t.Helper()

	if path != noMountPrefix {
		return PathJoin(mount, path)
	}
	return path
}

// compareErrors asserts each error in the list matches the provided error using errors.Is()
func compareErrors(t *testing.T, ev error, el []error) {
	t.Helper()

	for _, err := range el {
		assert.True(t, errors.Is(ev, err), fmt.Sprintf("error %v is not of type %v", ev, err))
	}
}

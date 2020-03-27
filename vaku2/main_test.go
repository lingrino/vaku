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

// noMountPrefix is a special string that, when passed in tests, will not be prefixed with a mount
// to allow testing on a nonexistent mount.
var noMountPrefix = "nomount"

// kvMountVersions lists the types of kv mounts for vault. There are currently two k/v mount types
// and vaku supports both. Tests should run against each version and return the same results.
var kvMountVersions = []string{"1", "2"}

// versionProduct is all possible to/from version combinations for testing functions that should
// work across multiple mount versions and vault servers.
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

// TestMain runs before any test. It is here because vault.TestCoreUnsealedWithConfig() calls a
// function further down that uses a default logger intead of the logger passed.
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

// testClient returns a client that points to a seeded server.
func testClient(t *testing.T, opts ...Option) (net.Listener, *Client) {
	t.Helper()

	ln, apiClient := testServerSeeded(t)

	client, err := NewClient(append(opts, WithVaultClient(apiClient))...)
	assert.NoError(t, err)

	return ln, client
}

// testClientDiffDst returns a client that points src and dst at different seeded servers.
func testClientDiffDst(t *testing.T, opts ...Option) (net.Listener, net.Listener, *Client) {
	t.Helper()

	ln, apiClientS := testServerSeeded(t)
	lnD, apiClientD := testServerSeeded(t)

	client, err := NewClient(append(opts,
		WithVaultSrcClient(apiClientS),
		WithVaultDstClient(apiClientD),
	)...)
	assert.NoError(t, err)

	return ln, lnD, client
}

// cloneCLient copies a client.
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
	op    string
	realL logical
}

// verify compliance with logical interface.
var _ logical = (*errLogical)(nil)

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

// updateLogical updates a client's real logical clients with passed errLogical clients.
func updateLogical(t *testing.T, c *Client, srcL logical, dstL logical) {
	t.Helper()

	if srcL != nil {
		sl, ok := srcL.(*errLogical)
		if ok {
			sl.realL = c.srcL
			c.srcL = sl
		} else {
			c.srcL = srcL
		}
	}
	if dstL != nil {
		dl, ok := dstL.(*errLogical)
		if ok {
			dl.realL = c.dstL
			c.dstL = dl
		} else {
			c.dstL = dstL
		}
	}
}

// addMountToPath prefixes a path with a mount if path is not noMountPrefix.
func addMountToPath(t *testing.T, path string, mount string) string {
	t.Helper()

	if path != noMountPrefix {
		return PathJoin(mount, path)
	}
	return path
}

// compareErrors asserts that the error list is an ordered and complete list of errors returned by
// repeatedly calling errors.Unwrap(err).
func compareErrors(t *testing.T, err error, el []error) {
	t.Helper()

	for _, e := range el {
		assert.True(t, errors.Is(err, e), fmt.Sprintf("error %v is not of type %v", err, e))
		err = errors.Unwrap(err)
	}

	assert.Nil(t, err)
}

func TestE(t *testing.T) {
	var errOne = errors.New("one")
	var errSecond = errors.New("second")
	var errTree = errors.New("tree")

	err := newWrapErr("1", errOne, nil)
	err = newWrapErr("2", errSecond, err)
	err = newWrapErr("3", errTree, err)

	exp := []error{errTree, errSecond, errOne}

	compareErrors(t, err, exp)
}

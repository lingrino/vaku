package vaku

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	vl "github.com/hashicorp/vault/sdk/logical"
)

// sharedVaku for most vaku tests. Tests isolate by path on each mount.
var sharedVaku *Client

// pathPrefix is used to create new unique seeded paths. Should be incremented after use.
var pathPrefix int
var pathPrefixMtx sync.Mutex

// mountless, when passed in tests, will not be prefixed with a mount.
const mountless = "mountless"

// mountVersions lists all kv versions. Tests run against all versions with equal results.
var mountVersions = [2]string{"1", "2"}

// versionProduct is all possible to/from mount version combinations.
var versionProduct = [4][2]string{
	{"1", "1"},
	{"2", "2"},
	{"1", "2"},
	{"2", "1"},
}

// seeds is the canonical secret seeds for every test.
var seeds = map[string]map[string]interface{}{
	"1": {
		"2": "3",
	},
	"4/5": {
		"6": "7",
	},
	"4/8": {
		"9":  "10",
		"11": "12",
	},
	"4/13/14": {
		"15": "16",
	},
	"4/13/17": {
		"18": "19",
		"20": "21",
		"22": "23",
	},
	"4/13/24/25/26/27": {
		"28": "29",
	},
}

// TestMain prepares the test run.
func TestMain(m *testing.M) {
	hclog.DefaultOutput = ioutil.Discard
	os.Exit(m.Run())
}

// testServer creates a new vault server and returns a vault API client that points to it.
func testServer(t *testing.T) *api.Client {
	t.Helper()

	// create vault core
	core, _, token := vault.TestCoreUnsealedWithConfig(t, &vault.CoreConfig{
		// Must be provided for v1/v2 path differences to work.
		LogicalBackends: map[string]vl.Factory{
			"kv": kv.Factory,
		},
		Logger: hclog.NewNullLogger(),
	})
	ln, addr := http.TestServer(t, core)
	t.Cleanup(func() { ln.Close() })

	// create client that points at core
	client, err := api.NewClient(api.DefaultConfig())
	assert.NoError(t, err)
	client.SetToken(token)
	assert.NoError(t, client.SetAddress(addr))

	// mount all mount versions
	for _, ver := range mountVersions {
		err := client.Sys().Mount(ver+"/", &api.MountInput{
			Type: "kv",
			Options: map[string]string{
				"version": ver,
			},
		})
		assert.NoError(t, err)
	}

	return client
}

// seededPath seeds a new prefixed path on the shared client. Returns the prefix to use.
func seededPath(t *testing.T) string {
	t.Helper()

	pathPrefixMtx.Lock()
	prefix := strconv.Itoa(pathPrefix)
	pathPrefix++
	pathPrefixMtx.Unlock()

	for _, ver := range mountVersions {
		for path, secret := range seeds {
			err := sharedVaku.PathWrite(PathJoin(ver, prefix, path), secret)
			assert.NoError(t, err)

			err = sharedVaku.dc.PathWrite(PathJoin(ver, prefix, path), secret)
			assert.NoError(t, err)
		}
	}

	return prefix
}

// testClient returns a client that points to a seeded server.
func testClient(t *testing.T, opts ...Option) *Client {
	t.Helper()

	apiClient := testServerSeeded(t)

	client, err := NewClient(append(opts, WithVaultClient(apiClient))...)
	assert.NoError(t, err)

	return client
}

// testClientDiffDst returns a client that points src and dst at different seeded servers.
func testClientDiffDst(t *testing.T, opts ...Option) *Client {
	t.Helper()

	apiClientS := testServerSeeded(t)
	apiClientD := testServerSeeded(t)

	client, err := NewClient(append(opts,
		WithVaultSrcClient(apiClientS),
		WithVaultDstClient(apiClientD),
	)...)
	assert.NoError(t, err)

	return client
}

// cloneCLient copies a client.
func cloneCLient(t *testing.T, c *Client) *Client {
	t.Helper()

	dc := *c.dc
	cpy := *c
	cpy.dc = &dc
	return &cpy
}

// testSetup sets up most of our tests. Returns a client with 'logical' updated and a readback client.
func testSetup(t *testing.T, srcL, dstL logical, opts ...Option) (*Client, *Client) {
	client := testClient(t, opts...)
	rbClient := cloneCLient(t, client)
	updateLogical(t, client, srcL, dstL)

	return client, rbClient
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
			sl.realL = c.vl
			c.vl = sl
		} else {
			c.vl = srcL
		}
	}
	if dstL != nil {
		dl, ok := dstL.(*errLogical)
		if ok {
			dl.realL = c.dc.vl
			c.dc.vl = dl
		} else {
			c.dc.vl = dstL
		}
	}
}

// addMountToPath prefixes a path with a mount if path is not mountless.
func addMountToPath(t *testing.T, path string, mount string) string {
	t.Helper()

	if path != mountless {
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

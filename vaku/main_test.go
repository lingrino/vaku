package vaku

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

// sharedReadBack for reading back in tests. No logical injections here.
var sharedReadBack *Client

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
	"0/1": {
		"2": "3",
	},
	"0/4/5": {
		"6": "7",
	},
	"0/4/8": {
		"9":  "10",
		"11": "12",
	},
	"0/4/13/14": {
		"15": "16",
	},
	"0/4/13/17": {
		"18": "19",
		"20": "21",
		"22": "23",
	},
	"0/4/13/24/25/26/27": {
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
	_, addr := http.TestServer(t, core)

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

// seededPath seeds a new prefixed path, appends to given path, and returns paths to test against.
func seededPath(t *testing.T, p string) []string {
	t.Helper()

	// set up the shared client if needed
	pathPrefixMtx.Lock()
	if sharedVaku == nil {
		initSharedVaku(t)
	}

	// use current path prefix and increment
	prefix := strconv.Itoa(pathPrefix)
	pathPrefix++
	pathPrefixMtx.Unlock()

	// mountless is for testing operations against invalid mounts
	if p == mountless {
		return []string{""}
	}

	// seed prefixes
	prefixes := make([]string, len(mountVersions))
	for i, ver := range mountVersions {
		for p, secret := range seeds {
			err := sharedVaku.PathWrite(PathJoin(ver, prefix, p), secret)
			assert.NoError(t, err)

			err = sharedVaku.dc.PathWrite(PathJoin(ver, prefix, p), secret)
			assert.NoError(t, err)
		}
		prefixes[i] = PathJoin(ver, prefix)
	}

	return prefixes
}

func initSharedVaku(t *testing.T) {
	t.Helper()

	srcClient := testServer(t)
	dstClient := testServer(t)

	client, err := NewClient(
		WithVaultSrcClient(srcClient),
		WithVaultDstClient(dstClient),
		WithabsolutePath(false),
		WithWorkers(100),
	)
	assert.NoError(t, err)

	// sharedReadBack does not have logical injector
	cleanDC := *client.dc
	cleanClient := *client
	cleanClient.dc = &cleanDC
	sharedReadBack = &cleanClient

	// replace standard logical with logicalInjector
	li := &logicalInjector{realL: client.vl}
	client.vl = li
	client.dc.vl = li
	sharedVaku = client
}

// inject is an injection to return from any logical function.
type inject struct {
	secret *api.Secret
	err    error
}

// injects is a map of path endings to an inject to return for that path.
var injects = map[string]inject{
	"injecterror":       {err: errInject},
	"injectdatanil":     {secret: &api.Secret{Data: nil}},
	"injectkeysnil":     {secret: &api.Secret{Data: map[string]interface{}{"keys": nil}}},
	"injectkeysint":     {secret: &api.Secret{Data: map[string]interface{}{"keys": 1}}},
	"injectkeyslistint": {secret: &api.Secret{Data: map[string]interface{}{"keys": []interface{}{1}}}},
}

// logicalInjector injects errors and outputs into vault operations.
type logicalInjector struct {
	realL logical
}

// verify compliance with logical interface.
var _ logical = (*logicalInjector)(nil)

func (e *logicalInjector) Delete(p string) (*api.Secret, error) {
	inject, ok := injects[path.Base(p)]
	if !ok {
		return e.realL.Delete(p)
	}
	return inject.secret, inject.err
}

func (e *logicalInjector) List(p string) (*api.Secret, error) {
	inject, ok := injects[path.Base(p)]
	if !ok {
		return e.realL.List(p)
	}
	return inject.secret, inject.err
}

func (e *logicalInjector) Read(p string) (*api.Secret, error) {
	inject, ok := injects[path.Base(p)]
	if !ok {
		return e.realL.Read(p)
	}
	return inject.secret, inject.err
}

func (e *logicalInjector) Write(p string, data map[string]interface{}) (*api.Secret, error) {
	i, ok := injects[path.Base(p)]
	if !ok {
		return e.realL.Write(p, data)
	}
	return i.secret, i.err
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

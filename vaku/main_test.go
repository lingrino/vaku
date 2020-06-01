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

// mountless, when passed in tests, will not be prefixed with a mount.
const mountless = "mountless"

// sharedVaku for most vaku tests. Tests isolate by path on each mount.
var sharedVaku *Client

// sharedVakuClean is for reading back values in tests. No logical injections here.
var sharedVakuClean *Client

// pathPrefix is used to create new unique seeded paths. Should be incremented after use.
var pathPrefix int = 100
var pathPrefixMtx sync.Mutex

// mountVersions lists all kv versions. Tests run against all versions with equal results.
var mountVersions = [2]string{"1", "2"}

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
		err := client.Sys().Mount("kv"+ver+"/", &api.MountInput{
			Type: "kv",
			Options: map[string]string{
				"version": ver,
			},
		})
		assert.NoError(t, err)
	}

	return client
}

// initSharedVaku sets up the global clients.
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

	// sharedVakuClean does not have logical injector
	cleanDC := *client.dc
	cleanClient := *client
	cleanClient.dc = &cleanDC
	cleanClient.vl = &logicalInjector{realL: client.vl, t: t, disabled: true}
	cleanClient.dc.vl = &logicalInjector{realL: client.dc.vl, t: t, disabled: true}
	sharedVakuClean = &cleanClient

	// replace standard logical with logicalInjector
	li := &logicalInjector{realL: client.vl, t: t}
	lid := &logicalInjector{realL: client.dc.vl, t: t}
	client.vl = li
	client.dc.vl = lid
	sharedVaku = client
}

// seededPrefixes seeds a new prefixed path, appends to given path, and returns paths to test against.
func seededPrefixes(t *testing.T, p string) []string {
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
			err := sharedVaku.PathWrite(PathJoin("kv"+ver, prefix, p), secret)
			assert.NoError(t, err)

			err = sharedVaku.dc.PathWrite(PathJoin("kv"+ver, prefix, p), secret)
			assert.NoError(t, err)
		}
		prefixes[i] = PathJoin("kv"+ver, prefix)
	}

	return prefixes
}

// seededPrefixProduct returns a list of pairs of prefixes to use for src/dest commands.
func seededPrefixProduct(t *testing.T) [4][2]string {
	t.Helper()

	prefixes1 := seededPrefixes(t, "")
	prefixes2 := seededPrefixes(t, "")
	prefixes3 := seededPrefixes(t, "")
	prefixes4 := seededPrefixes(t, "")

	return [4][2]string{
		{prefixes1[0], prefixes1[0]},
		{prefixes2[0], prefixes2[1]},
		{prefixes3[1], prefixes3[0]},
		{prefixes4[1], prefixes4[1]},
	}
}

// testName takes a path (and optional destination path) and returns a name for the test.
func testName(sp string, dps ...string) string {
	var dp string
	if len(dps) == 1 {
		dp = dps[0]
	}

	if dp != "" {
		return fmt.Sprintf("~%s->%s~", sp, dp)
	}

	return fmt.Sprintf("~%s~", sp)
}

func TestComareErrors(t *testing.T) {
	e := newWrapErr("", ErrPathWrite, newWrapErr("", ErrVaultWrite, nil))
	ee := []error{ErrPathWrite, ErrVaultWrite}
	compareErrors(t, e, ee)
}

// compareErrors asserts continuously calling Unwrap(err) produces the error list.
func compareErrors(t *testing.T, err error, el []error) {
	t.Helper()

	for _, e := range el {
		assert.True(t, errors.Is(err, e), fmt.Sprintf("error %v is not of type %v", err, e))
		err = errors.Unwrap(err)
	}

	assert.Nil(t, err)
}

// inject represents an injection to return instead of a real vault API call.
type inject struct {
	secret *api.Secret
	err    error
}

// injects is a map of path endings to an inject to return for that path.
var injects = map[string]*inject{
	"error":       {err: errInject},
	"nildata":     {secret: &api.Secret{Data: nil}},
	"nilkeys":     {secret: &api.Secret{Data: map[string]interface{}{"keys": nil}}},
	"intkeys":     {secret: &api.Secret{Data: map[string]interface{}{"keys": 1}}},
	"listintkeys": {secret: &api.Secret{Data: map[string]interface{}{"keys": []interface{}{1}}}},
	"funcdata": {secret: &api.Secret{Data: map[string]interface{}{
		"data": map[string]interface{}{
			"foo": func() {},
		},
		"metadata": map[string]interface{}{
			"destroyed":     false,
			"deletion_time": "",
		},
	}}},
}

// logicalInjector injects errors and outputs into vault api calls.
type logicalInjector struct {
	t        *testing.T
	realL    logical
	disabled bool
}

// verify compliance with logical interface.
var _ logical = (*logicalInjector)(nil)

// run does injector logic. inject at a path with normalpath/injectname/operation/inject.
func (e *logicalInjector) run(p, op string) (string, *inject) {
	e.t.Helper()

	// if not injecting, proceed as normal
	if path.Base(p) != "inject" {
		return p, nil
	}
	p = path.Dir(p)

	// if not injecting on this operation, proceed with dir path
	if path.Base(p) != op {
		return path.Dir(path.Dir(p)), nil
	}
	p = path.Dir(p)

	// if no injector or disabled, proceed with dir path
	inj, ok := injects[path.Base(p)]
	if !ok || e.disabled {
		return path.Dir(p), nil
	}

	// return injector
	return path.Dir(p), inj
}

func (e *logicalInjector) Delete(p string) (*api.Secret, error) {
	e.t.Helper()

	p, inj := e.run(p, "delete")
	if inj != nil {
		return inj.secret, inj.err
	}
	return e.realL.Delete(p)
}

func (e *logicalInjector) List(p string) (*api.Secret, error) {
	e.t.Helper()

	p, inj := e.run(p, "list")
	if inj != nil {
		return inj.secret, inj.err
	}
	return e.realL.List(p)
}

func (e *logicalInjector) Read(p string) (*api.Secret, error) {
	e.t.Helper()

	p, inj := e.run(p, "read")
	if inj != nil {
		return inj.secret, inj.err
	}
	return e.realL.Read(p)
}

func (e *logicalInjector) Write(p string, data map[string]interface{}) (*api.Secret, error) {
	e.t.Helper()

	p, inj := e.run(p, "write")
	if inj != nil {
		return inj.secret, inj.err
	}
	return e.realL.Write(p, data)
}

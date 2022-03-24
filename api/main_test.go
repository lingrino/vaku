package vaku

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
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
var pathPrefix = 100
var pathPrefixMtx sync.Mutex

// mountVersions lists all kv versions. Tests run against all versions with equal results.
var mountVersions = [2]string{"1", "2"}

// seeds is the canonical secret seeds for every test.
var seeds = map[string]map[string]any{
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
	hclog.DefaultOutput = io.Discard
	os.Exit(m.Run())
}

// testServer creates a new vault server and returns a vault API client that points to it.
// Pass an empty &testing.T{} to vt if you're initializing a long-lived client so that
// vt.Cleanup() does not shutdown your shared client.
func testServer(t *testing.T, vt *testing.T) *api.Client {
	t.Helper()

	// create vault core
	core, _, token := vault.TestCoreUnsealedWithConfig(vt, &vault.CoreConfig{
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

	srcClient := testServer(t, &testing.T{})
	dstClient := testServer(t, &testing.T{})

	client, err := NewClient(
		WithVaultSrcClient(srcClient),
		WithVaultDstClient(dstClient),
		WithAbsolutePath(false),

		// set worker < max folder operations to expose any worker threading issue
		WithWorkers(5),
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

// seededPrefixes seeds a new prefixed path, appends to given path, and returns paths to test
// against.
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
		seedsCopy := make(map[string]map[string]any, len(seeds))
		for p, v := range seeds {
			seedsCopy[PathJoin("kv"+ver, prefix, p)] = v
		}

		err := sharedVakuClean.FolderWrite(context.Background(), seedsCopy)
		assert.NoError(t, err)

		err = sharedVakuClean.dc.FolderWrite(context.Background(), seedsCopy)
		assert.NoError(t, err)

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
	"nilkeys":     {secret: &api.Secret{Data: map[string]any{"keys": nil}}},
	"intkeys":     {secret: &api.Secret{Data: map[string]any{"keys": 1}}},
	"listintkeys": {secret: &api.Secret{Data: map[string]any{"keys": []any{1}}}},
	"funcdata": {secret: &api.Secret{Data: map[string]any{
		"data": map[string]any{
			"foo": func() {},
		},
		"metadata": map[string]any{
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

// run does injector logic. inject at a path with path/injectname/operation/inject/path.
func (e *logicalInjector) run(p, op string) (string, *inject) {
	e.t.Helper()

	// remove trailing slash
	p = strings.TrimSuffix(p, "/")

	// find injection in path
	var injectOp string
	var injectName string
	pathSplit := strings.Split(p, "/")
	for i, s := range pathSplit {
		if s == "inject" {
			injectOp = pathSplit[i-1]
			injectName = pathSplit[i-2]
			pathSplit[i] = ""
			pathSplit[i-1] = ""
			pathSplit[i-2] = ""
		}
	}

	// if not injecting, proceed as normal
	if injectName == "" {
		return p, nil
	}

	// cleanPath is path with injection words removed
	cleanPath := PathJoin(pathSplit...)

	// if not injecting on this operation, proceed with dir path
	if injectOp != op {
		return cleanPath, nil
	}

	// if no injector or disabled, proceed with dir path
	inj, ok := injects[injectName]
	if !ok || e.disabled {
		return cleanPath, nil
	}

	// return injector
	return cleanPath, inj
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

	np, inj := e.run(p, "list")

	// return injected results if they exist
	if inj != nil {
		return inj.secret, inj.err
	}

	// call list with the real logical client
	sec, errL := e.realL.List(np)

	// if we are injecting but not on list, then we need to re-add the injections to the path of
	// the results so that the injections can continue with the resulting paths.
	if path.Base(p) == "inject" && !e.disabled {
		if sec == nil || sec.Data == nil {
			return nil, errL
		}
		list, err := decodeSecret(sec)
		if err != nil {
			return nil, err
		}

		listI := []any{}
		for _, l := range list {
			p = strings.TrimSuffix(p, "/")
			ip := PathJoin(path.Base(path.Dir(path.Dir(p))), path.Base(path.Dir(p)), path.Base(p))
			np = PathJoin(l, ip)
			if IsFolder(l) {
				np = EnsureFolder(np)
			}
			listI = append(listI, np)
		}

		ns := &api.Secret{
			Data: map[string]any{
				"keys": listI,
			},
		}
		return ns, errL
	}

	return sec, errL
}

func (e *logicalInjector) Read(p string) (*api.Secret, error) {
	e.t.Helper()

	p, inj := e.run(p, "read")
	if inj != nil {
		return inj.secret, inj.err
	}
	return e.realL.Read(p)
}

func (e *logicalInjector) Write(p string, data map[string]any) (*api.Secret, error) {
	e.t.Helper()

	p, inj := e.run(p, "write")
	if inj != nil {
		return inj.secret, inj.err
	}
	return e.realL.Write(p, data)
}

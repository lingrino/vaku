package cmd

import (
	"bytes"
	"context"
	"testing"

	"github.com/lingrino/vaku/vaku"
)

// newTestCLI returns a CLI ready for running tests.
func newTestCLI(t *testing.T, args []string) (*cli, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	cli := newCLI()
	cli.flagIndent = ""

	var outW, errW bytes.Buffer
	cli.cmd.SetOut(&outW)
	cli.cmd.SetErr(&errW)

	cli.cmd.SetArgs(args)

	return cli, &outW, &errW
}

// newTestCLIWithAPI returns a CLI with an initialized API ready for running tests.
func newTestCLIWithAPI(t *testing.T, args []string) (*cli, *bytes.Buffer, *bytes.Buffer) {
	t.Helper()

	cli, outW, errW := newTestCLI(t, args)
	cli.vc = &testVakuClient{}
	return cli, outW, errW
}

// testVakuClient implements vaku.ClientInterface and just returns very basic values. We don't need
// to test the vaku client again, it is already well tested.
type testVakuClient struct{}

// Verify Client compliance with the interface.
var _ vaku.ClientInterface = (*testVakuClient)(nil)

func (c *testVakuClient) FolderCopy(ctx context.Context, src, dst string) error {
	return nil
}

func (c *testVakuClient) FolderDelete(ctx context.Context, p string) error {
	return nil
}

func (c *testVakuClient) FolderList(ctx context.Context, p string) ([]string, error) {
	return []string{"foo/bar", "foo/baz", "bim/bom"}, nil
}

func (c *testVakuClient) FolderListChan(ctx context.Context, p string) (<-chan string, <-chan error) {
	return nil, nil
}

func (c *testVakuClient) FolderMove(ctx context.Context, src, dst string) error {
	return nil
}

func (c *testVakuClient) FolderRead(ctx context.Context, p string) (map[string]map[string]interface{}, error) {
	return map[string]map[string]interface{}{
		"foo": {
			"bim": "bom",
			"biz": "baz",
		},
		"bar": {
			"hoo": "boo",
		},
	}, nil
}

func (c *testVakuClient) FolderReadChan(ctx context.Context, p string) (<-chan map[string]map[string]interface{}, <-chan error) { //nolint:lll
	return nil, nil
}

func (c *testVakuClient) FolderSearch(ctx context.Context, path, search string) ([]string, error) {
	return []string{"foo/bar", "bim/bom"}, nil
}

func (c *testVakuClient) FolderWrite(ctx context.Context, d map[string]map[string]interface{}) error {
	return nil
}

func (c *testVakuClient) PathCopy(src, dst string) error {
	return nil
}

func (c *testVakuClient) PathDelete(p string) error {
	return nil
}

func (c *testVakuClient) PathList(p string) ([]string, error) {
	return []string{"foo", "moo"}, nil
}

func (c *testVakuClient) PathMove(src, dst string) error {
	return nil
}

func (c *testVakuClient) PathRead(p string) (map[string]interface{}, error) {
	return map[string]interface{}{"biz": "baz", "foo": "bar"}, nil
}

func (c *testVakuClient) PathSearch(p, s string) (bool, error) {
	return true, nil
}

func (c *testVakuClient) PathUpdate(p string, d map[string]interface{}) error {
	return nil
}

func (c *testVakuClient) PathWrite(p string, d map[string]interface{}) error {
	return nil
}

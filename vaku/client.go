package vaku

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

var (
	ErrNumWorkers = errors.New("invalid workers")
)

const (
	defaultWorkers = 10
)

// logical is functions from api.Logical() used by Vaku. Helps with testing.
type logical interface {
	Delete(path string) (*api.Secret, error)
	List(path string) (*api.Secret, error)
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
}

// Client has all Vaku functions and wraps Vault API clients.
type Client struct {
	// vc is the vault client.
	vc *api.Client
	// vl wraps vc.Logical() for easy testing.
	vl logical

	// dc is a recursive Client for operations with a source and destination.
	dc *Client

	// workers is the max number of concurrent operations against vault.
	workers int

	// absolutePath if the absolute path is desired instead of the relative path.
	absolutePath bool
}

// Option configures a Client.
type Option interface {
	apply(c *Client) error
}

// WithVaultClient sets the Vault client to be used.
func WithVaultClient(c *api.Client) Option {
	return withVaultClient{c}
}

// WithVaultSrcClient is an alias for WithVaultClient.
func WithVaultSrcClient(c *api.Client) Option {
	return withVaultClient{c}
}

type withVaultClient struct {
	client *api.Client
}

func (o withVaultClient) apply(c *Client) error {
	c.vc = o.client
	c.vl = o.client.Logical()
	return nil
}

// WithVaultDstClient sets a separate Vault client to be used only on operations that have a source
// and destination (copy, move, etc...). If unset the source client will be used.
func WithVaultDstClient(c *api.Client) Option {
	return withVaultDstClient{c}
}

type withVaultDstClient struct {
	client *api.Client
}

func (o withVaultDstClient) apply(c *Client) error {
	c.dc.vc = o.client
	c.dc.vl = o.client.Logical()
	return nil
}

// WithWorkers sets the maximum number of goroutines that access Vault at any given time. Does not
// cap the number of goroutines overall. Default value is 10. A stable and well-operated Vault
// server should be able to handle 100 or more without issue. Use with caution and tune specifically
// to your environment and storage backend.
func WithWorkers(n int) Option {
	return withWorkers(n)
}

type withWorkers int

func (o withWorkers) apply(c *Client) error {
	if o < 1 {
		return newWrapErr(fmt.Sprintf("workers must 1 or greater: %d", o), ErrNumWorkers, nil)
	}
	c.workers = int(o)
	c.dc.workers = int(o)
	return nil
}

// WithabsolutePath sets the output format for all returned paths. Default path output is a relative
// path, trimmed up to the path input. Pass WithabsolutePath(true) to set path output to the entire
// path. Example: List(secret/foo) -> "bar" OR "secret/foo/bar".
func WithabsolutePath(b bool) Option {
	return withabsolutePath(b)
}

type withabsolutePath bool

func (o withabsolutePath) apply(c *Client) error {
	c.absolutePath = bool(o)
	c.dc.absolutePath = bool(o)
	return nil
}

// NewClient returns a new Vaku Client based on the Vault API config.
func NewClient(opts ...Option) (*Client, error) {
	// set defaults
	client := &Client{
		workers: defaultWorkers,
	}
	client.dc = client

	// apply options
	for _, opt := range opts {
		err := opt.apply(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

// pathToRetur returns a path for the user, given their formatting preferences.
func (c *Client) pathToReturn(path, root string) string {
	if c.absolutePath {
		return EnsurePrefix(path, root)
	}
	return PathJoin(strings.TrimPrefix(path, root))
}

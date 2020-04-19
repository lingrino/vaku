package vaku

import (
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	defaultWorkers = 10
)

// logical is functions from api.Logical() used by Vaku.
type logical interface {
	Delete(path string) (*api.Secret, error)
	List(path string) (*api.Secret, error)
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
}

// Client has all Vaku functions and wraps Vault API clients.
type Client struct {
	vc *api.Client
	vl logical
	dc *Client

	// src is the default client and also used as dst when dst is nil.
	src *api.Client
	dst *api.Client

	// wrap api.Client.Logical() in an interface.
	srcL logical
	dstL logical

	// workers is the max number of concurrent operations against vault.
	workers int

	// absolutePath if the absolute path is desired instead of the relative path.
	absolutePath bool
}

// Option configures a Client.
type Option interface {
	apply(c *Client) error
}

// WithVaultClient sets the default Vault client to be used.
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
	c.src = o.client
	c.srcL = o.client.Logical()
	return nil
}

// WithVaultDstClient sets a separate Vault client to be used only on operations that have a source
// and destination (copy, move, etc...). If unset the default client will be used as the source and
// destination.
func WithVaultDstClient(c *api.Client) Option {
	return withDstVaultClient{c}
}

type withDstVaultClient struct {
	client *api.Client
}

func (o withDstVaultClient) apply(c *Client) error {
	c.dc = &Client{
		vc: o.client,
		vl: o.client.Logical(),
		dc: c,
	}
	c.dst = o.client
	c.dstL = o.client.Logical()
	return nil
}

// WithWorkers sets the maximum number of goroutines that access Vault at any given time. Does not
// cap the number of goroutines overall. Default value is 10. A stable and well-operated Vault
// server should be able to handle 100 or more without issue. Use with caution and tune specifically
// to your environment and storage backend.
func WithWorkers(n int) Option {
	return withWorkers(n)
}

type withWorkers uint

func (o withWorkers) apply(c *Client) error {
	c.workers = int(o)
	return nil
}

// WithabsolutePath sets the output format for all returned paths. Default path output is a relative
// path, trimmed up to the path input. Pass WithabsolutePath(true) to set path output to the entire
// path. Example: List(secret/foo) -> "bar" OR "secret/foo/bar"
func WithabsolutePath(b bool) Option {
	return withabsolutePath(b)
}

type withabsolutePath bool

func (o withabsolutePath) apply(c *Client) error {
	c.absolutePath = bool(o)
	return nil
}

// NewClient returns a new Vaku Client based on the Vault API config.
func NewClient(opts ...Option) (*Client, error) {
	// set defaults
	client := &Client{
		workers: defaultWorkers,
	}

	// apply options
	for _, opt := range opts {
		err := opt.apply(client)
		if err != nil {
			return nil, err
		}
	}

	// set dst to src if dst is unspecified
	if client.dst == nil && client.src != nil {
		client.dst = client.src
		client.dstL = client.srcL
	}

	if client.dc == nil {
		client.dc = client
	}

	return client, nil
}

// pathToReturn takes a path and returns a path that can be returned to the user, given their
// formatting preferences.
func (c *Client) pathToReturn(path, root string) string {
	if c.absolutePath {
		return EnsurePrefix(path, root)
	}
	return KeyJoin(strings.TrimPrefix(path, root))
}

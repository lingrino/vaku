package vaku2

import (
	"github.com/hashicorp/vault/api"
)

// logical has all functions from api.Logical() that I use
type logical interface {
	Delete(path string) (*api.Secret, error)
	List(path string) (*api.Secret, error)
	Read(path string) (*api.Secret, error)
	Write(path string, data map[string]interface{}) (*api.Secret, error)
}

// Client holds Vaku functions and wraps Vault API clients.
type Client struct {
	// source is the default client and also used as dest when dest is nil.
	source *api.Client
	dest   *api.Client

	// wrap api.Client.Logical() in an interface for mocking
	sourceL logical
	destL   logical

	// max number of concurrent operations we'll run.
	workers uint

	// set for the full path to be returned instead of the trimmed path
	fullPath bool
}

// Option configures a Client.
type Option interface {
	apply(c *Client) error
}

type withVaultClient struct {
	client *api.Client
}

func (o withVaultClient) apply(c *Client) error {
	c.source = o.client
	c.sourceL = o.client.Logical()
	return nil
}

// WithVaultClient sets the default vault client to be used
func WithVaultClient(c *api.Client) Option {
	return withVaultClient{c}
}

// WithVaultSourceClient is an alias for WithVaultClient
func WithVaultSourceClient(c *api.Client) Option {
	return withVaultClient{c}
}

type withDestVaultClient struct {
	client *api.Client
}

func (o withDestVaultClient) apply(c *Client) error {
	c.dest = o.client
	c.destL = o.client.Logical()
	return nil
}

// WithVaultDestClient sets a separate vault client to be used only on operations that have a source
// and destination (copy, move, etc...). If unset the default client will be used as the source and
// destination.
func WithVaultDestClient(c *api.Client) Option {
	return withDestVaultClient{c}
}

type withWorkers uint

func (o withWorkers) apply(c *Client) error {
	c.workers = uint(o)
	return nil
}

// WithWorkers sets the maximum number of goroutines that will be used to run folder based
// functions. The default value is 10, but a stable and well-tuned Vault server should be able to
// handle up to 100 without issues. Use with caution and tune specifically to your environment and
// storage backend.
func WithWorkers(n uint) Option {
	return withWorkers(n)
}

type withFullPath bool

func (o withFullPath) apply(c *Client) error {
	c.fullPath = bool(o)
	return nil
}

// WithFullPath sets the output format for all returned paths. By default path output is trimmed up
// to the path input. Pass WithFullPath(true) to set path output to the entire path. Example:
// List(secret/foo) -> "bar" OR "secret/foo/bar"
func WithFullPath(b bool) Option {
	return withFullPath(b)
}

// NewClient returns a new empty Vaku Client based on the Vault API config
func NewClient(opts ...Option) (*Client, error) {
	// set defaults
	client := &Client{
		workers: 10,
	}

	// apply options
	for _, opt := range opts {
		err := opt.apply(client)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

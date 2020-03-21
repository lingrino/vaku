package vaku2

import (
	"github.com/hashicorp/vault/api"
)

// Client holds Vaku functions and wraps Vault API clients.
type Client struct {
	// vault clients used to call the API
	// source is also the client used when there is no dest set
	source *api.Client
	dest   *api.Client

	// number of concurrent operations we'll run at once
	workers uint
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

func WithWorkers(n uint) Option {
	return withWorkers(n)
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

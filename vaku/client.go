package vaku

import (
	vapi "github.com/hashicorp/vault/api"
)

// Client is a wrapper around a real Vault API client.
type Client struct {
	*vapi.Client
}

// NewClient Returns a new empty Client type
func NewClient() *Client {
	return &Client{}
}

package vaku

import (
	vapi "github.com/hashicorp/vault/api"
)

// Client is a simple wrapper around a real Vault API client. All Vaku
// functions are defined on Client as well, which lets anyone already
// using a Vault Client easily make use of Vaku.
type Client struct {
	*vapi.Client
}

// NewClient returns a new empty Vaku Client. Using this function requires
// that you initialize and set the nested Vault client on your own before
// Vaku functions will work.
func NewClient() *Client {
	return &Client{}
}

// NewClientFromVaultClient takes in an official Vault Client that you have
// already created and returns a Vaku client that wraps your Vault client.
func NewClientFromVaultClient(vc *vapi.Client) *Client {
	vakuClient := NewClient()
	vakuClient.Client = vc
	return vakuClient
}

// CopyClient is a client for copy operations where the source address/namespace/token is different from
// target address/namespace/token. The source is a client for the source of the copy, target is a client for the target
// of the copy
type CopyClient struct {
	Source *Client
	Target *Client
}

// NewCopyClient returns a new empty CopyClient.  Using this function requires
// that caller initialize and set the source / target client.
func NewCopyClient() *CopyClient {
	return &CopyClient{}
}

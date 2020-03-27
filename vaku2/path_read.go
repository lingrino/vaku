package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrVaultRead = errors.New("vault read")
)

// PathRead takes a path, calls vault read, extracts the secret, and returns it.
func (c *Client) PathRead(p string) (map[string]interface{}, error) {
	return c.pathRead(c.sourceL, p)
}

// PathReadDest takes a path, calls vault read using the destination client, extracts the secret, and returns it.
func (c *Client) PathReadDest(p string) (map[string]interface{}, error) {
	return c.pathRead(c.destL, p)
}

// pathRead takes a path, calls vault read, extracts the secret, and returns it.
func (c *Client) pathRead(apiL logical, p string) (map[string]interface{}, error) {
	secret, err := apiL.Read(p)
	if err != nil {
		return nil, newWrapErr(fmt.Sprintf("%q: %v: %v", p, ErrVaultRead, err), ErrVaultRead, nil)
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	return secret.Data, nil
}

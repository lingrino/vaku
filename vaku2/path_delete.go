package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrVaultDelete = errors.New("vault delete")
)

// PathDelete deletes data at a path using the source client
func (c *Client) PathDelete(p string) error {
	return c.pathDelete(c.sourceL, p)
}

// PathDeleteDest deletes data at a path using the dest client
func (c *Client) PathDeleteDest(p string) error {
	return c.pathDelete(c.destL, p)
}

// pathDelete deletes data at a path
func (c *Client) pathDelete(apiL logical, p string) error {
	_, err := apiL.Delete(p)
	if err != nil {
		return fmt.Errorf("%q: %w", p, ErrVaultDelete)
	}

	return nil
}

package vaku

import (
	"errors"
)

var (
	// ErrPathDelete when PathDelete/PathDeleteDst fails.
	ErrPathDelete = errors.New("path delete")
	// ErrVaultDelete when the underlying Vault API delete fails.
	ErrVaultDelete = errors.New("vault delete")
)

// PathDelete deletes data at a path.
func (c *Client) PathDelete(p string) error {
	return c.pathDelete(c.srcL, p)
}

// PathDeleteDst deletes data at a path.
func (c *Client) PathDeleteDst(p string) error {
	return c.pathDelete(c.dstL, p)
}

// pathDelete does the actual delete.
func (c *Client) pathDelete(l logical, p string) error {
	_, err := l.Delete(p)
	if err != nil {
		return newWrapErr(p, ErrPathDelete, newWrapErr(err.Error(), ErrVaultDelete, nil))
	}

	return nil
}

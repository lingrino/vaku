package vaku

import (
	"errors"
)

var (
	// ErrPathDelete when PathDelete fails.
	ErrPathDelete = errors.New("path delete")
	// ErrVaultDelete when the underlying Vault API delete fails.
	ErrVaultDelete = errors.New("vault delete")
)

// PathDelete deletes data at a path.
func (c *Client) PathDelete(p string) error {
	_, err := c.vl.Delete(p)
	if err != nil {
		return newWrapErr(p, ErrPathDelete, newWrapErr(err.Error(), ErrVaultDelete, nil))
	}

	return nil
}

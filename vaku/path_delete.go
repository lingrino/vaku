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
	vaultPath, _, err := c.rewritePath(p, vaultDelete)
	if err != nil {
		return newWrapErr(p, ErrPathDelete, err)
	}

	_, err = c.vl.Delete(vaultPath)
	if err != nil {
		return newWrapErr(p, ErrPathDelete, newWrapErr(err.Error(), ErrVaultDelete, nil))
	}

	return nil
}

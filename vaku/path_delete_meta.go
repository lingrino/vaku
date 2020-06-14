package vaku

import (
	"errors"
)

var (
	// ErrPathDeleteMeta when PathDeleteMeta fails.
	ErrPathDeleteMeta = errors.New("path delete meta")
)

// PathDeleteMeta deletes all secret metadata and versions. Only works on v2 kv engines.
func (c *Client) PathDeleteMeta(p string) error {
	vaultPath, _, err := c.rewritePath(p, vaultDeleteMeta)
	if err != nil {
		return newWrapErr(p, ErrPathDeleteMeta, err)
	}

	_, err = c.vl.Delete(vaultPath)
	if err != nil {
		return newWrapErr(p, ErrPathDeleteMeta, newWrapErr(err.Error(), ErrVaultDelete, nil))
	}

	return nil
}

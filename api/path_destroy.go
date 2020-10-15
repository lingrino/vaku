package vaku

import (
	"errors"
)

var (
	// ErrPathDestroy when PathDestroy fails.
	ErrPathDestroy = errors.New("path destroy")
)

// PathDestroy destroys versions of a secret at a path. Only works on v2 kv engines.
func (c *Client) PathDestroy(p string, versions []int) error {
	if len(versions) == 0 {
		return newWrapErr("no versions provided", ErrPathDestroy, nil)
	}

	vaultPath, _, err := c.rewritePath(p, vaultDestroy)
	if err != nil {
		return newWrapErr(p, ErrPathDestroy, err)
	}

	data := map[string]interface{}{
		"versions": versions,
	}

	_, err = c.vl.Write(vaultPath, data)
	if err != nil {
		return newWrapErr(p, ErrPathDestroy, newWrapErr(err.Error(), ErrVaultWrite, nil))
	}

	return nil
}

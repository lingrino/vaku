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
	err := c.pathDeleteWithOp(p, vaultDeleteMeta)
	if err != nil {
		return newWrapErr(p, ErrPathDeleteMeta, err)
	}

	return err
}

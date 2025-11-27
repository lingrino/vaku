package vaku

import (
	"errors"
)

var (
	// ErrPathMoveAllVersions when PathMoveAllVersions fails.
	ErrPathMoveAllVersions = errors.New("path move all versions")
)

// PathMoveAllVersions moves all versions of a secret from source to destination (copy + delete).
// Only works on v2 kv engines.
func (c *Client) PathMoveAllVersions(src, dst string) error {
	err := c.PathCopyAllVersions(src, dst)
	if err != nil {
		return newWrapErr("", ErrPathMoveAllVersions, err)
	}

	err = c.PathDeleteMeta(src)
	if err != nil {
		return newWrapErr(src, ErrPathMoveAllVersions, err)
	}

	return nil
}

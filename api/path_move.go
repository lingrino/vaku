package vaku

import (
	"errors"
)

var (
	// ErrPathMove when PathMove fails.
	ErrPathMove = errors.New("path move")
)

// PathMove moves data at a source path to a destination path (copy + delete).
// If allVersions is true, all versions of the secret are moved (KV v2 only).
func (c *Client) PathMove(src, dst string, allVersions bool) error {
	err := c.PathCopy(src, dst, allVersions)
	if err != nil {
		return newWrapErr("", ErrPathMove, err)
	}

	err = c.PathDelete(src)
	if err != nil {
		return newWrapErr(dst, ErrPathMove, err)
	}

	return nil
}

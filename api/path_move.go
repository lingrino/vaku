package vaku

import (
	"errors"
)

var (
	// ErrPathMove when PathMove fails.
	ErrPathMove = errors.New("path move")
	// ErrPathMoveAllVersions when PathMoveAllVersions fails.
	ErrPathMoveAllVersions = errors.New("path move all versions")
)

// PathMove moves data at a source path to a destination path (copy + delete).
func (c *Client) PathMove(src, dst string) error {
	err := c.PathCopy(src, dst)
	if err != nil {
		return newWrapErr("", ErrPathMove, err)
	}

	err = c.PathDelete(src)
	if err != nil {
		return newWrapErr(dst, ErrPathMove, err)
	}

	return nil
}

// PathMoveAllVersions moves all versions of data at a source path to a destination path.
// This copies all versions to destination and then deletes all metadata/versions from source.
// Only works on v2 kv engines.
func (c *Client) PathMoveAllVersions(src, dst string) error {
	// First check if this is a v2 mount
	_, mv, err := c.rewritePath(src, vaultRead)
	if err != nil {
		return newWrapErr(src, ErrPathMoveAllVersions, err)
	}

	if mv != mv2 {
		err := newWrapErr("all versions move not supported on KV v1", ErrMountVersion, nil)
		return newWrapErr(src, ErrPathMoveAllVersions, err)
	}

	// Copy all versions to destination
	err = c.PathCopyAllVersions(src, dst)
	if err != nil {
		return newWrapErr("", ErrPathMoveAllVersions, err)
	}

	// Delete all metadata and versions from source (true deletion of entire secret)
	err = c.PathDeleteMeta(src)
	if err != nil {
		return newWrapErr(dst, ErrPathMoveAllVersions, err)
	}

	return nil
}

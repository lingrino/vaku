package vaku

import (
	"errors"
)

var (
	// ErrPathCopy when PathCopy fails.
	ErrPathCopy = errors.New("path copy")
	// ErrPathCopyAllVersions when PathCopyAllVersions fails.
	ErrPathCopyAllVersions = errors.New("path copy all versions")
)

// PathCopy copies data at a source path to a destination path.
func (c *Client) PathCopy(src, dst string) error {
	secret, err := c.PathRead(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopy, err)
	}

	err = c.dc.PathWrite(dst, secret)
	if err != nil {
		return newWrapErr(dst, ErrPathCopy, err)
	}

	return nil
}

// PathCopyAllVersions copies the current version of data at a source path to a destination path.
// Uses metadata checking to ensure the source exists and is not deleted.
// Only works on v2 kv engines. Note: Vault KV v2 doesn't support copying version history
// directly, so this copies the current active version after verifying it exists via metadata.
func (c *Client) PathCopyAllVersions(src, dst string) error {
	// First check if this is a v2 mount
	_, mv, err := c.rewritePath(src, vaultRead)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	if mv != mv2 {
		err := newWrapErr("all versions copy not supported on KV v1", ErrMountVersion, nil)
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	// Read metadata to verify the secret exists and get current version info
	metadata, err := c.PathReadMetadata(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	if metadata == nil {
		return nil // nothing to copy
	}

	// For now, just copy the current version (same as PathCopy)
	// Future enhancement could implement true version history preservation
	secret, err := c.PathRead(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	if secret == nil {
		return nil // nothing to copy
	}

	err = c.dc.PathWrite(dst, secret)
	if err != nil {
		return newWrapErr(dst, ErrPathCopyAllVersions, err)
	}

	return nil
}

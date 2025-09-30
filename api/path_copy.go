package vaku

import (
	"errors"
)

var (
	// ErrPathCopy when PathCopy fails.
	ErrPathCopy = errors.New("path copy")
)

// PathCopy copies data at a source path to a destination path.
func (c *Client) PathCopy(src, dst string) error {
	return c.pathCopy(src, dst, false)
}

// PathCopyAllVersions copies all versions of a secret at a source path to a destination path (KV v2 only).
func (c *Client) PathCopyAllVersions(src, dst string) error {
	return c.pathCopy(src, dst, true)
}

// pathCopy is the internal implementation for path copying.
func (c *Client) pathCopy(src, dst string, allVersions bool) error {
	if !allVersions {
		// Original behavior: copy only the latest version
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

	// Copy all versions
	versions, err := c.PathReadAllVersions(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopy, err)
	}

	// Write each version in order
	for _, version := range versions {
		err = c.dc.PathWrite(dst, version)
		if err != nil {
			return newWrapErr(dst, ErrPathCopy, err)
		}
	}

	return nil
}

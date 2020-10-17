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

package vaku2

import (
	"errors"
)

var (
	ErrPathCopy = errors.New("path copy")
)

// PathCopy copies data at a source path to a destination path. Client must have been initialized
// using WithDstClient() when copying across vault servers.
func (c *Client) PathCopy(src, dst string) error {
	secret, err := c.PathRead(src)
	if err != nil {
		return newWrapErr("read "+src, ErrPathCopy, err)
	}

	err = c.PathWriteDst(dst, secret)
	if err != nil {
		return newWrapErr("write "+dst, ErrPathCopy, err)
	}

	return nil
}

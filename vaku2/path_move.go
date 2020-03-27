package vaku2

import (
	"errors"
)

var (
	ErrPathMove = errors.New("path move")
)

// PathMove moves data at a source path to a destination path (copy + delete). Client must
// have been initialized using WithDstClient() when moving across vault servers.
func (c *Client) PathMove(src, dst string) error {
	err := c.PathCopy(src, dst)
	if err != nil {
		return newWrapErr("copy "+src+" to "+dst, ErrPathMove, err)
	}

	err = c.PathDelete(src)
	if err != nil {
		return newWrapErr("delete "+dst, ErrPathMove, err)
	}

	return nil
}

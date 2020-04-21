package vaku

import (
	"context"
	"errors"
)

var (
	// ErrFolderMove when FolderMove fails.
	ErrFolderMove = errors.New("folder move")
)

// FolderMove moves data at a source folder to a destination folder. Source is deleted after copy.
// Client must have been initialized using WithDstClient() when moving across vault servers.
func (c *Client) FolderMove(ctx context.Context, src, dst string) error {
	err := c.dc.FolderCopy(ctx, src, dst)
	if err != nil {
		return newWrapErr("", ErrFolderMove, err)
	}

	err = c.dc.FolderDelete(ctx, src)
	if err != nil {
		return newWrapErr("delete "+src, ErrFolderMove, err)
	}

	return nil
}

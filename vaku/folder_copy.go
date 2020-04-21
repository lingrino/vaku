package vaku

import (
	"context"
	"errors"
)

var (
	// ErrFolderCopy when FolderCopy fails.
	ErrFolderCopy = errors.New("folder copy")
)

// FolderCopy copies data at a source folder to a destination folder. Client must have been
// initialized using WithDstClient() when copying across vault servers.
func (c *Client) FolderCopy(ctx context.Context, src, dst string) error {
	read, err := c.FolderRead(ctx, src)
	if err != nil {
		return newWrapErr("read from "+src, ErrFolderCopy, err)
	}

	// Switch the key prefixes from src to dst
	if c.absolutePath {
		TrimMapKeyPrefix(read, src)
	}
	EnsureMapKeyPrefix(read, dst)

	err = c.dc.FolderWrite(ctx, read)
	if err != nil {
		return newWrapErr("write to "+dst, ErrFolderCopy, err)
	}

	return nil
}

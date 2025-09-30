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
func (c *Client) FolderMove(ctx context.Context, src, dst string) error {
	err := c.FolderCopy(ctx, src, dst)
	if err != nil {
		return newWrapErr("", ErrFolderMove, err)
	}

	err = c.FolderDelete(ctx, src)
	if err != nil {
		return newWrapErr("delete "+src, ErrFolderMove, err)
	}

	return nil
}

// FolderMoveAllVersions moves all versions of each secret at a source folder to a destination folder (KV v2 only). Source is deleted after copy.
func (c *Client) FolderMoveAllVersions(ctx context.Context, src, dst string) error {
	err := c.FolderCopyAllVersions(ctx, src, dst)
	if err != nil {
		return newWrapErr("", ErrFolderMove, err)
	}

	err = c.FolderDelete(ctx, src)
	if err != nil {
		return newWrapErr("delete "+src, ErrFolderMove, err)
	}

	return nil
}

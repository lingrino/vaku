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
// If allVersions is true, all versions of each secret are moved (KV v2 only).
func (c *Client) FolderMove(ctx context.Context, src, dst string, allVersions bool) error {
	err := c.FolderCopy(ctx, src, dst, allVersions)
	if err != nil {
		return newWrapErr("", ErrFolderMove, err)
	}

	err = c.FolderDelete(ctx, src)
	if err != nil {
		return newWrapErr("delete "+src, ErrFolderMove, err)
	}

	return nil
}

package vaku

import (
	"context"
	"errors"
)

var (
	// ErrFolderMoveAllVersions when FolderMoveAllVersions fails.
	ErrFolderMoveAllVersions = errors.New("folder move all versions")
)

// FolderMoveAllVersions moves all versions of all secrets from source folder to destination folder.
// Only works on v2 kv engines for both source and destination.
func (c *Client) FolderMoveAllVersions(ctx context.Context, src, dst string) error {
	err := c.FolderCopyAllVersions(ctx, src, dst)
	if err != nil {
		return newWrapErr("", ErrFolderMoveAllVersions, err)
	}

	err = c.FolderDeleteMeta(ctx, src)
	if err != nil {
		return newWrapErr(src, ErrFolderMoveAllVersions, err)
	}

	return nil
}

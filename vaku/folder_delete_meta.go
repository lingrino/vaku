package vaku

import (
	"context"
	"errors"
)

var (
	// ErrFolderDeleteMeta when FolderDeleteMeta fails.
	ErrFolderDeleteMeta = errors.New("folder delete meta")
)

// FolderDeleteMeta deletes all secret metadata and versions for all secrets in a folder. Only works
// on v2 kv engines.
func (c *Client) FolderDeleteMeta(ctx context.Context, p string) error {
	return nil
}

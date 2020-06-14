package vaku

import (
	"context"
	"errors"
)

var (
	// ErrFolderDestroy when FolderDestroy fails.
	ErrFolderDestroy = errors.New("folder destroy")
)

// FolderDestroy destroys versions of all secrets in a folder. Only works on v2 kv engines.
func (c *Client) FolderDestroy(ctx context.Context, p string, versions []int) error {
	return nil
}

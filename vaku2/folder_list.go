package vaku2

import (
	"errors"
)

var (
	ErrFolderList = errors.New("folder list")
)

// FolderList lists the provided path and all subpaths.
func (c *Client) FolderList(p string) ([]string, error) {
	return nil, nil
}

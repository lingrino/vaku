package vaku

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrFolderSearch when FolderSearch fails.
	ErrFolderSearch = errors.New("folder search")
)

// FolderSearch searches the provided path and all subpaths. Returns a list of
// paths in which the string was found.
func (c *Client) FolderSearch(ctx context.Context, path, search string) ([]string, error) {
	read, err := c.FolderRead(ctx, path)
	if err != nil {
		return nil, newWrapErr(path, ErrFolderSearch, err)
	}

	var matches []string
	for pth, sec := range read {
		fmt.Println(sec)
		found, err := searchSecret(sec, search)
		if err != nil {
			return nil, newWrapErr(path, ErrFolderSearch, err)
		}
		if found {
			matches = append(matches, pth)
		}
	}

	return matches, nil
}

package vaku

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderDelete when FolderDelete fails.
	ErrFolderDelete = errors.New("folder delete")
	// ErrFolderDeleteChan when FolderDeleteChan fails.
	ErrFolderDeleteChan = errors.New("folder delete chan")
)

// FolderDelete recursively deletes the provided path and all subpaths.
func (c *Client) FolderDelete(ctx context.Context, p string) error {
	errC := c.FolderDeleteChan(ctx, p)

	err := <-errC
	if err != nil {
		return newWrapErr(p, ErrFolderDelete, err)
	}

	return nil
}

// FolderDeleteChan recursively deletes the provided path and all subpaths. Returns an error channel
// that sends either the first error or nil when the work is done.
func (c *Client) FolderDeleteChan(ctx context.Context, p string) <-chan error {
	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)

	// list the path
	pathC, errC := c.FolderListChan(ctx, p)
	eg.Go(func() error {
		err := <-errC
		if err != nil {
			return newWrapErr(p, ErrFolderDeleteChan, err)
		}
		return nil
	})

	// fan out and process paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error {
			return c.folderDeleteWork(&folderDeleteWorkInput{
				ctx:   ctx,
				root:  p,
				pathC: pathC,
			})
		})
	}

	return errFuncOnChan(eg.Wait)
}

// folderDeleteWorkInput is the piecces needed to list a folder
type folderDeleteWorkInput struct {
	ctx   context.Context
	root  string
	pathC <-chan string
}

// folderDeleteWork takes input from pathC, lists the path, adds listed folders back into pathC, and
// adds non-folders into results.
func (c *Client) folderDeleteWork(i *folderDeleteWorkInput) error {
	for {
		select {
		case <-i.ctx.Done():
			return i.ctx.Err()
		case path, ok := <-i.pathC:
			if !ok {
				return nil
			}
			path = EnsurePrefix(path, i.root)
			err := c.PathDelete(path)
			if err != nil {
				return newWrapErr(i.root, ErrFolderDeleteChan, err)
			}
			return nil
		}
	}
}

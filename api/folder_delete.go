package vaku

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderDelete when FolderDelete fails.
	ErrFolderDelete = errors.New("folder delete")
)

// FolderDelete recursively deletes the provided path and all subpaths.
func (c *Client) FolderDelete(ctx context.Context, p string) error {
	err := c.folderDeleteWithFunc(ctx, p, c.PathDelete)
	if err != nil {
		return newWrapErr(p, ErrFolderDelete, err)
	}

	return nil
}

func (c *Client) folderDeleteWithFunc(ctx context.Context, p string, deleteF func(string) error) error {
	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)

	// list the path
	pathC, errC := c.FolderListChan(ctx, p)
	eg.Go(func() error {
		err := <-errC
		if err != nil {
			return err
		}
		return nil
	})

	// fan out and process paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error {
			return c.folderDeleteWork(&folderDeleteWorkInput{
				ctx:     ctx,
				root:    p,
				pathC:   pathC,
				deleteF: deleteF,
			})
		})
	}

	return eg.Wait() //nolint:wrapcheck
}

// folderDeleteWorkInput is the piecces needed to list a folder.
type folderDeleteWorkInput struct {
	ctx     context.Context
	root    string
	pathC   <-chan string
	deleteF func(string) error
}

// folderDeleteWork takes input from pathC, lists the path, adds listed folders back into pathC, and
// adds non-folders into results.
func (c *Client) folderDeleteWork(i *folderDeleteWorkInput) error {
	for {
		select {
		case <-i.ctx.Done():
			return ctxErr(i.ctx.Err())
		case path, ok := <-i.pathC:
			if !ok {
				return nil
			}
			path = EnsurePrefix(path, i.root)
			err := i.deleteF(path)
			if err != nil {
				return err
			}
			return nil
		}
	}
}

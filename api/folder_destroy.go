package vaku

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderDestroy when FolderDestroy fails.
	ErrFolderDestroy = errors.New("folder destroy")
)

// FolderDestroy destroys versions of all secrets in a folder. Only works on v2 kv engines.
func (c *Client) FolderDestroy(ctx context.Context, p string, versions []int) error {
	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)

	// list the path
	pathC, errC := c.FolderListChan(ctx, p)
	eg.Go(func() error {
		return <-errC
	})

	// fan out and process paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error {
			return c.folderDestroyWork(&folderDestroyWorkInput{
				ctx:      ctx,
				root:     p,
				versions: versions,
				pathC:    pathC,
			})
		})
	}

	err := eg.Wait()
	if err != nil {
		return newWrapErr(p, ErrFolderDestroy, err)
	}
	return nil
}

// folderDestroyWorkInput is the piecces needed to destroy a folder.
type folderDestroyWorkInput struct {
	ctx      context.Context
	root     string
	versions []int
	pathC    <-chan string
}

// folderDestroyWork takes input from pathC, lists the path, adds listed folders back into pathC, and
// adds non-folders into results.
func (c *Client) folderDestroyWork(i *folderDestroyWorkInput) error {
	for {
		select {
		case <-i.ctx.Done():
			return ctxErr(i.ctx.Err())
		case path, ok := <-i.pathC:
			if !ok {
				return nil
			}
			path = EnsurePrefix(path, i.root)
			return c.PathDestroy(path, i.versions)
		}
	}
}

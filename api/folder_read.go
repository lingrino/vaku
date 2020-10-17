package vaku

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderRead when FolderRead fails.
	ErrFolderRead = errors.New("folder read")
	// ErrFolderReadChan when FolderReadChan fails.
	ErrFolderReadChan = errors.New("folder read chan")
)

// FolderRead recursively reads the provided path and all subpaths.
func (c *Client) FolderRead(ctx context.Context, p string) (map[string]map[string]interface{}, error) {
	resC, errC := c.FolderReadChan(ctx, p)

	// read results and errors. send on errC signifies done (can be nil).
	out := make(map[string]map[string]interface{})
	for {
		select {
		case res := <-resC:
			mergeMaps(out, res)
		case err := <-errC:
			if err != nil {
				return nil, newWrapErr(p, ErrFolderRead, err)
			}
			if len(out) == 0 {
				return nil, nil
			}
			return out, nil
		}
	}
}

// FolderReadChan recursively reads the provided path and all subpaths. Returns an unbuffered
// channel that can be read until close and an error channel that sends either the first error or
// nil when the work is done.
func (c *Client) FolderReadChan(ctx context.Context, p string) (<-chan map[string]map[string]interface{}, <-chan error) { //nolint:lll
	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)

	// resC is processed paths
	resC := make(chan map[string]map[string]interface{})

	// list the path
	pathC, errC := c.FolderListChan(ctx, p)
	eg.Go(func() error {
		err := <-errC
		if err != nil {
			return newWrapErr(p, ErrFolderReadChan, err)
		}
		return nil
	})

	// fan out and process paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error {
			return c.folderReadWork(&folderReadWorkInput{
				ctx:   ctx,
				root:  p,
				pathC: pathC,
				resC:  resC,
			})
		})
	}

	return resC, errFuncOnChan(eg.Wait)
}

// folderReadWorkInput is the piecces needed to list a folder.
type folderReadWorkInput struct {
	ctx   context.Context
	root  string
	pathC <-chan string
	resC  chan<- map[string]map[string]interface{}
}

// folderReadWork takes input from pathC, lists the path, adds listed folders back into pathC, and
// adds non-folders into results.
func (c *Client) folderReadWork(i *folderReadWorkInput) error {
	for {
		select {
		case <-i.ctx.Done():
			return i.ctx.Err()
		case path, ok := <-i.pathC:
			if !ok {
				return nil
			}
			err := c.pathReadWork(path, i)
			if err != nil {
				return err
			}
		}
	}
}

// pathReadWork reads the path adds results to the channel.
func (c *Client) pathReadWork(path string, i *folderReadWorkInput) error {
	path = EnsurePrefix(path, i.root)

	read, err := c.PathRead(path)
	if err != nil {
		return newWrapErr(i.root, ErrFolderReadChan, err)
	}

	// Don't add nil reads to results. These show up in list but are actually deleted secrets.
	if read != nil {
		res := make(map[string]map[string]interface{}, 1)
		res[c.outputPath(path, i.root)] = read

		i.resC <- res
	}

	return nil
}

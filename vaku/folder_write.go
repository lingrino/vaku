package vaku

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderWrite when FolderWrite fails.
	ErrFolderWrite = errors.New("folder write")
	// ErrFolderWriteChan when FolderWriteChan fails.
	ErrFolderWriteChan = errors.New("folder write chan")
)

// // FolderWrite writes data to a path. Multiple paths can be written to at once.
// func (c *Client) FolderWrite(ctx context.Context, d map[string]map[string]interface{}) error {
// 	return c.folderWrite(ctx, d)
// }

// FolderWrite writes data to a path. Multiple paths can be written to at once.
func (c *Client) FolderWrite(ctx context.Context, d map[string]map[string]interface{}) error {
	errC := c.FolderWriteChan(ctx, d)

	err := <-errC
	if err != nil {
		return newWrapErr("", ErrFolderWrite, err)
	}

	return nil
}

// FolderWriteChan writes data to a path. Multiple paths can be written to at once. Returns an error channel
// that sends either the first error or nil when the work is done.
func (c *Client) FolderWriteChan(ctx context.Context, d map[string]map[string]interface{}) <-chan error {
	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)

	// add paths to be processed by our workers
	pathC := make(chan string)
	go func() {
		for path := range d {
			pathC <- path
		}
		close(pathC)
	}()

	// fan out and process paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error {
			return c.folderWriteWork(&folderWriteWorkInput{
				ctx:   ctx,
				pathC: pathC,
				data:  d,
			})
		})
	}

	return errFuncOnChan(eg.Wait)
}

// folderWriteWorkInput is the piecces needed to list a folder
type folderWriteWorkInput struct {
	ctx   context.Context
	pathC <-chan string
	data  map[string]map[string]interface{}
}

// folderWriteWork takes input from pathC, lists the path, adds listed folders back into pathC, and
// adds non-folders into results.
func (c *Client) folderWriteWork(i *folderWriteWorkInput) error {
	for {
		select {
		case <-i.ctx.Done():
			return i.ctx.Err()
		case path, ok := <-i.pathC:
			if !ok {
				return nil
			}
			err := c.PathWrite(path, i.data[path])
			if err != nil {
				return newWrapErr(path, ErrFolderWriteChan, err)
			}

			return nil
		}
	}
}

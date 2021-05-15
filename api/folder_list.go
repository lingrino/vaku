package vaku

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderList when FolderList fails.
	ErrFolderList = errors.New("folder list")
	// ErrFolderListChan when FolderListChan fails.
	ErrFolderListChan = errors.New("folder list chan")
)

// FolderList recursively lists the provided path and all subpaths.
func (c *Client) FolderList(ctx context.Context, p string) ([]string, error) {
	resC, errC := c.FolderListChan(ctx, p)

	// read results and errors. send on errC signifies done (can be nil).
	var output []string
	for {
		select {
		case res, ok := <-resC:
			if !ok {
				return output, nil
			}
			output = append(output, res)
		case err := <-errC:
			if err != nil {
				return nil, newWrapErr(p, ErrFolderList, err)
			}
		}
	}
}

// FolderListChan recursively lists the provided path and all subpaths. Returns an unbuffered
// channel that can be read until close and an error channel that sends either the first error or
// nil when the work is done.
func (c *Client) FolderListChan(ctx context.Context, p string) (<-chan string, <-chan error) {
	// input must be a folder (end in "/")
	root := EnsureFolder(p)

	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)
	// wg tracks when to close the paths channel
	var wg sync.WaitGroup

	// pathC is paths to be processed
	pathC := make(chan string)
	// resC is processed paths
	resC := make(chan string)
	// errC for the first error seen
	errC := make(chan error)

	// add root path to paths
	wg.Add(1)
	go func(p string) { pathC <- p }(root)

	// fan out and process paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error {
			return c.folderListWork(&folderListWorkInput{
				ctx:   ctx,
				root:  root,
				wg:    &wg,
				pathC: pathC,
				resC:  resC,
			})
		})
	}

	// Wait until finished (success or not) and clean up
	go func() {
		// Close pathC after all paths added
		wg.Wait()
		close(pathC)

		// Wait for all paths to process
		err := eg.Wait()

		// Report the error (or nil) to errC
		errC <- err

		// Clean up
		close(resC)
		close(errC)
	}()

	return resC, errC
}

// folderListWorkInput is the pieces needed to list a folder.
type folderListWorkInput struct {
	ctx   context.Context
	root  string
	wg    *sync.WaitGroup
	pathC chan string
	resC  chan<- string
}

// folderListWork takes input from pathC, lists the path, adds listed folders back into pathC, and
// adds non-folders into results.
func (c *Client) folderListWork(i *folderListWorkInput) error {
	for path := range i.pathC {
		err := c.pathListWork(path, i)
		if err != nil {
			return err
		}

		select {
		case <-i.ctx.Done():
			return ctxErr(i.ctx.Err())
		default:
		}
	}
	return nil
}

// pathListWork takes a path and either adds it back to the pathC (if folder) or processes it and
// adds it to the resC.
func (c *Client) pathListWork(path string, i *folderListWorkInput) error {
	defer i.wg.Done()

	if IsFolder(path) {
		list, err := c.PathList(path)
		if err != nil {
			return newWrapErr(i.root, ErrFolderListChan, err)
		}
		for _, item := range list {
			i.wg.Add(1)
			go func(item string) {
				item = EnsurePrefix(item, path)
				i.pathC <- item
			}(item)
		}
	} else {
		select {
		case i.resC <- c.outputPath(path, i.root):
		case <-i.ctx.Done():
			return ctxErr(i.ctx.Err())
		}
	}
	return nil
}

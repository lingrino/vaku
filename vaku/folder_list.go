package vaku

import (
	"context"
	"errors"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderList when FolderList fails.
	ErrFolderList = errors.New("folder list")
)

// FolderList recursively lists the provided path and all subpaths.
func (c *Client) FolderList(ctx context.Context, p string) ([]string, error) {
	resC, errC := c.FolderListChan(ctx, p)

	// read results and errors. send on errC signifies done (can be nil).
	var out []string
	for {
		select {
		case res := <-resC:
			out = append(out, res)
		case err := <-errC:
			if err != nil {
				return nil, err
			}
			return out, nil
		}
	}
}

// FolderListChan recursively lists the provided path and all subpaths. Returns an unbuffered
// channel that can be read until close and an error channel that sends either the first error or
// nil when the work is done.
func (c *Client) FolderListChan(ctx context.Context, p string) (<-chan string, <-chan error) {
	// input must be a folder (end in "/")
	root := MakeFolder(p)

	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)
	// wg tracks when to close the paths channel
	var wg sync.WaitGroup

	// pathC is paths to be processed
	pathC := make(chan string)
	// resC is processed paths
	resC := make(chan string)

	// add root path to paths
	wg.Add(1)
	go func(p string) { pathC <- p }(root)

	// fan out and process paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error { return c.folderListWork(ctx, root, &wg, pathC, resC) })
	}

	// close pathC once all have been processed or when the group is cancelled
	eg.Go(func() error {
		// provide a way to wait on wg.Wait() inside a select
		done := make(chan bool)
		go func() {
			wg.Wait()
			done <- true
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-done:
			close(pathC)
			return nil
		}
	})

	// provide eg.Wait() on a channel for returning
	errC := make(chan error)
	go func() { errC <- eg.Wait() }()

	return resC, errC
}

// folderListWork takes input from pathC, lists the path, adds listed folders back into pathC, and
// adds non-folders into results.
func (c *Client) folderListWork(ctx context.Context, root string, wg *sync.WaitGroup, pathC, resC chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case path, ok := <-pathC:
			if !ok {
				return nil
			}
			return c.pathListWork(path, root, wg, pathC, resC)
		}
	}
}

// pathListWork takes a path and either adds it back to the pathC (if folder) or processes it and
// adds it to the resC.
func (c *Client) pathListWork(path, root string, wg *sync.WaitGroup, pathC, resC chan string) error {
	if IsFolder(path) {
		list, err := c.PathList(path)
		if err != nil {
			return newWrapErr(root, ErrFolderList, err)
		}
		for _, item := range list {
			item = EnsurePrefix(item, path)
			wg.Add(1)
			go func(p string) { pathC <- p }(item)
		}
	} else {
		if c.absolutepath {
			resC <- path
		} else {
			resC <- strings.TrimPrefix(path, root)
		}
	}
	wg.Done()
	return nil
}

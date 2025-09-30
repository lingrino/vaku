package vaku

import (
	"context"
	"errors"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderCopy when FolderCopy fails.
	ErrFolderCopy = errors.New("folder copy")
)

// FolderCopy copies data at a source folder to a destination folder.
// If allVersions is true, all versions of each secret are copied (KV v2 only).
func (c *Client) FolderCopy(ctx context.Context, src, dst string, allVersions bool) error {
	if !allVersions {
		// Original behavior: copy only latest versions
		read, err := c.FolderRead(ctx, src)
		if err != nil {
			return newWrapErr("read from "+src, ErrFolderCopy, err)
		}

		// Switch the key prefixes from src to dst
		c.swapPaths(read, src, dst)

		err = c.dc.FolderWrite(ctx, read)
		if err != nil {
			return newWrapErr("write to "+dst, ErrFolderCopy, err)
		}

		return nil
	}

	// Copy all versions: list all paths and copy each with allVersions=true
	paths, err := c.FolderList(ctx, src)
	if err != nil {
		return newWrapErr("list "+src, ErrFolderCopy, err)
	}

	// eg manages workers processing paths
	eg, ctx := errgroup.WithContext(ctx)

	// add paths to be processed by our workers
	pathC := make(chan string, len(paths))
	for _, path := range paths {
		pathC <- path
	}
	close(pathC)

	// fan out and copy paths
	for i := 0; i < c.workers; i++ {
		eg.Go(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctxErr(ctx.Err())
				case path, ok := <-pathC:
					if !ok {
						return nil
					}
					srcPath := c.inputPath(path, src)
					dstPath := c.dc.inputPath(path, dst)
					err := c.PathCopy(srcPath, dstPath, true)
					if err != nil {
						return newWrapErr("copy "+srcPath+" to "+dstPath, ErrFolderCopy, err)
					}
				}
			}
		})
	}

	err = eg.Wait()
	if err != nil {
		return err
	}

	return nil
}

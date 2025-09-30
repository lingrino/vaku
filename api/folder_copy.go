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
func (c *Client) FolderCopy(ctx context.Context, src, dst string) error {
	return c.folderCopy(ctx, src, dst, false)
}

// FolderCopyAllVersions copies all versions of each secret at a source folder to a destination folder (KV v2 only).
func (c *Client) FolderCopyAllVersions(ctx context.Context, src, dst string) error {
	return c.folderCopy(ctx, src, dst, true)
}

// folderCopy is the internal implementation for folder copying.
func (c *Client) folderCopy(ctx context.Context, src, dst string, allVersions bool) error {
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
					err := c.PathCopyAllVersions(srcPath, dstPath)
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

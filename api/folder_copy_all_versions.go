package vaku

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/sync/errgroup"
)

var (
	// ErrFolderCopyAllVersions when FolderCopyAllVersions fails.
	ErrFolderCopyAllVersions = errors.New("folder copy all versions")
)

// FolderCopyAllVersions copies all versions of all secrets from source folder to destination folder.
// Only works on v2 kv engines for both source and destination.
func (c *Client) FolderCopyAllVersions(ctx context.Context, src, dst string) error {
	// Validate both source and destination are KV v2 mounts
	if err := c.validateKV2Mount(src, c); err != nil {
		return newWrapErr(src, ErrFolderCopyAllVersions, err)
	}
	if err := c.validateKV2Mount(dst, c.dc); err != nil {
		return newWrapErr(dst, ErrFolderCopyAllVersions, err)
	}

	// eg manages workers reading from the paths channel
	eg, ctx := errgroup.WithContext(ctx)

	// list the path
	pathC, errC := c.FolderListChan(ctx, src)
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
			return c.folderCopyAllVersionsWork(&folderCopyAllVersionsWorkInput{
				ctx:   ctx,
				src:   src,
				dst:   dst,
				pathC: pathC,
			})
		})
	}

	err := eg.Wait()
	if err != nil {
		return newWrapErr(src, ErrFolderCopyAllVersions, err)
	}
	return nil
}

// folderCopyAllVersionsWorkInput is the pieces needed to copy all versions of a folder.
type folderCopyAllVersionsWorkInput struct {
	ctx   context.Context
	src   string
	dst   string
	pathC <-chan string
}

// folderCopyAllVersionsWork processes paths from pathC and copies all versions of each secret.
func (c *Client) folderCopyAllVersionsWork(i *folderCopyAllVersionsWorkInput) error {
	for {
		select {
		case <-i.ctx.Done():
			return ctxErr(i.ctx.Err())
		case path, ok := <-i.pathC:
			if !ok {
				return nil
			}
			srcPath := c.inputPath(path, i.src)
			// Transform path from source to destination:
			// - When absolutePath=true: path has source prefix, strip it and add dest prefix
			// - When absolutePath=false: path is relative, just add dest prefix
			var dstPath string
			if c.absolutePath {
				relativePath := strings.TrimPrefix(path, i.src)
				dstPath = PathJoin(i.dst, relativePath)
			} else {
				dstPath = c.inputPath(path, i.dst)
			}
			err := c.PathCopyAllVersions(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
}

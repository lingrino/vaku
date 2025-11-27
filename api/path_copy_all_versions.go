package vaku

import (
	"errors"
	"sort"
)

var (
	// ErrPathCopyAllVersions when PathCopyAllVersions fails.
	ErrPathCopyAllVersions = errors.New("path copy all versions")
)

// PathCopyAllVersions copies all versions of a secret from source to destination.
// Only works on v2 kv engines for both source and destination.
// Deleted or destroyed versions are preserved by writing empty secrets.
func (c *Client) PathCopyAllVersions(src, dst string) error {
	// Validate source is KV v2
	_, srcVersion, err := c.mountInfo(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}
	if srcVersion != mv2 {
		return newWrapErr(src, ErrPathCopyAllVersions, ErrMountVersion)
	}

	// Validate destination is KV v2
	_, dstVersion, err := c.dc.mountInfo(dst)
	if err != nil {
		return newWrapErr(dst, ErrPathCopyAllVersions, err)
	}
	if dstVersion != mv2 {
		return newWrapErr(dst, ErrPathCopyAllVersions, ErrMountVersion)
	}

	// Read metadata to get all versions
	meta, err := c.PathReadMeta(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	// If no metadata or no versions, nothing to copy
	if meta == nil || len(meta.Versions) == 0 {
		return nil
	}

	// Get sorted list of version numbers
	versions := make([]int, 0, len(meta.Versions))
	for v := range meta.Versions {
		versions = append(versions, v)
	}
	sort.Ints(versions)

	// Copy each version in order
	for _, v := range versions {
		vMeta := meta.Versions[v]

		var data map[string]any
		if vMeta.Deleted || vMeta.Destroyed {
			// Write empty secret to preserve version position
			data = map[string]any{}
		} else {
			// Read the version data
			data, err = c.PathReadVersion(src, v)
			if err != nil {
				return newWrapErr(src, ErrPathCopyAllVersions, err)
			}
			// If data is nil (shouldn't happen for non-deleted), use empty map
			if data == nil {
				data = map[string]any{}
			}
		}

		// Write to destination
		err = c.dc.PathWrite(dst, data)
		if err != nil {
			return newWrapErr(dst, ErrPathCopyAllVersions, err)
		}
	}

	return nil
}

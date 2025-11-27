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
	if err := c.validateKV2Mount(src, c); err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}
	if err := c.validateKV2Mount(dst, c.dc); err != nil {
		return newWrapErr(dst, ErrPathCopyAllVersions, err)
	}

	meta, err := c.PathReadMeta(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}
	if meta == nil || len(meta.Versions) == 0 {
		return nil
	}

	versions := sortedVersions(meta.Versions)
	for _, v := range versions {
		if err := c.copyVersion(src, dst, v, meta.Versions[v]); err != nil {
			return err
		}
	}

	return nil
}

// validateKV2Mount checks that a path is on a KV v2 mount.
func (c *Client) validateKV2Mount(path string, client *Client) error {
	_, version, err := client.mountInfo(path)
	if err != nil {
		return err
	}
	if version != mv2 {
		return ErrMountVersion
	}
	return nil
}

// sortedVersions returns version numbers sorted in ascending order.
func sortedVersions(versions map[int]SecretVersionMeta) []int {
	result := make([]int, 0, len(versions))
	for v := range versions {
		result = append(result, v)
	}
	sort.Ints(result)
	return result
}

// copyVersion copies a single version from source to destination.
func (c *Client) copyVersion(src, dst string, version int, vMeta SecretVersionMeta) error {
	data, err := c.getVersionData(src, version, vMeta)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	if err := c.dc.PathWrite(dst, data); err != nil {
		return newWrapErr(dst, ErrPathCopyAllVersions, err)
	}
	return nil
}

// getVersionData retrieves data for a version, returning empty map for deleted/destroyed versions.
func (c *Client) getVersionData(src string, version int, vMeta SecretVersionMeta) (map[string]any, error) {
	if vMeta.Deleted || vMeta.Destroyed {
		return map[string]any{}, nil
	}

	data, err := c.PathReadVersion(src, version)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return map[string]any{}, nil
	}
	return data, nil
}

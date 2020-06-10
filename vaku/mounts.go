package vaku

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrMountInfo when failing to get mount information about a path.
	ErrMountInfo = errors.New("mount info")
	// ErrListMounts when failing to list vault mounts.
	ErrListMounts = errors.New("list mounts")
	// ErrNoMount when path cannot be matched to a mount.
	ErrNoMount = errors.New("no matching mount")
	// ErrRewritePath when failing to rewrite the path with mount data.
	ErrRewritePath = errors.New("rewriting path")
)

// mountVersion represents possible vault kv mount versions.
type mountVersion int

const (
	mv0 mountVersion = iota
	mv1
	mv2
)

// vaultOperation represents a call made to vault (read, write, delete, etc...).
type vaultOperation int

// operations that are used in vaku.
const (
	vaultList vaultOperation = iota + 1
	vaultRead
	vaultWrite
	vaultDelete
	vaultDestroy
)

// mountInfo takes a path and returns the mount path and version.
func (c *Client) mountInfo(p string) (string, mountVersion, error) {
	mounts, err := c.vc.Sys().ListMounts()
	if err != nil {
		return "", mv0, newWrapErr(p, ErrMountInfo, newWrapErr(err.Error(), ErrListMounts, nil))
	}

	for mount, data := range mounts {
		// Ensure '/' so that no match on foo/bar/ when actual path is foo/barbar/
		mount = EnsureFolder(mount)
		if strings.HasPrefix(p, mount) {
			version, ok := data.Options["version"]
			if !ok {
				return mount, mv0, nil
			}

			return mount, mountStringToVersion(version), nil
		}
	}

	return "", mv0, newWrapErr(p, ErrMountInfo, newWrapErr(p, ErrNoMount, nil))
}

// mountStringToVersion converts a mount version string from vault to a MountVersion.
func mountStringToVersion(v string) mountVersion {
	version, err := strconv.Atoi(v)
	if err != nil {
		return mv0
	}

	return mountVersion(version)
}

// rewritePath rewrites a vault input path based on the mount version and operation.
func (c *Client) rewritePath(p string, op vaultOperation) (string, mountVersion, error) {
	mount, version, err := c.mountInfo(p)
	if err != nil {
		return "", mv0, newWrapErr(p, ErrRewritePath, err)
	}

	// only rewrite mv2 mounts
	if version != mv2 {
		return p, version, nil
	}

	switch op {
	case vaultList, vaultDestroy:
		p = InsertIntoPath(p, mount, "metadata")
	case vaultRead, vaultWrite, vaultDelete:
		p = InsertIntoPath(p, mount, "data")
	}

	return p, version, nil
}

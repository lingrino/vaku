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
	// ErrMountVersion when an operation is not supported on the mount version.
	ErrMountVersion = errors.New("mount version does not support operation")
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
	vaultDeleteMeta
	vaultReadMeta
)

// mountInfo takes a path and returns the mount path and version.
func (c *Client) mountInfo(p string) (string, mountVersion, error) {
	mounts, err := c.mountProvider.ListMounts()
	if err != nil {
		return "", mv0, newWrapErr(p, ErrMountInfo, newWrapErr(err.Error(), ErrListMounts, nil))
	}

	for _, mount := range mounts {
		// Ensure '/' so that no match on foo/bar/ when actual path is foo/barbar/
		mount.Path = EnsureFolder(mount.Path)
		if strings.HasPrefix(p, mount.Path) {
			if mount.Version == "" {
				return mount.Path, mv0, nil
			}

			return mount.Path, mountStringToVersion(mount.Version), nil
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

	// check if the operation is supported
	if !mountSupportsOperation(op, version) {
		return "", version, newWrapErr(p, ErrMountVersion, nil)
	}

	// only rewrite mv2 mounts
	if version != mv2 {
		return p, version, nil
	}

	switch op {
	case vaultList, vaultDeleteMeta, vaultReadMeta:
		p = InsertIntoPath(p, mount, "metadata")
	case vaultRead, vaultWrite, vaultDelete:
		p = InsertIntoPath(p, mount, "data")
	case vaultDestroy:
		p = InsertIntoPath(p, mount, "destroy")
	}

	return p, version, nil
}

func mountSupportsOperation(op vaultOperation, v mountVersion) bool {
	// v2 mounts support all operations
	if v == mv2 {
		return true
	}

	// v1 mounts don't support these
	switch op {
	case vaultDestroy, vaultDeleteMeta, vaultReadMeta:
		return false
	}

	// default to supported
	return true
}

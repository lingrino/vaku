package vaku

import (
	"path"
	"strings"
)

// KeyIsFolder returns true if the string ends in a '/'.
func (c *Client) KeyIsFolder(k string) bool {
	if strings.HasSuffix(k, "/") {
		return true
	}
	return false
}

// KeyJoin takes n strings and combines them into a clean Vault key. Vault keys
// (as opposed to paths) optionally have a trailing '/' to signify they are a folder.
func (c *Client) KeyJoin(ks ...string) string {
	if strings.HasSuffix(ks[len(ks)-1], "/") {
		return strings.TrimPrefix(path.Join(ks...)+"/", "/")
	}
	return strings.TrimPrefix(path.Join(ks...), "/")
}

// PathJoin takes n strings and combines them into a clean Vault path. Paths in
// vault (unlike keys) never end with a '/'.
func (c *Client) PathJoin(ps ...string) string {
	return strings.TrimSuffix(c.KeyJoin(ps...), "/")
}

// KeyClean simply calls KeyJoin() with one string. KeyJoin already cleans keys, so
// this helper exists only for naming simplicity.
func (c *Client) KeyClean(k string) string {
	return c.KeyJoin(k)
}

// PathClean just calls PathJoin() with one string. PathJoin already cleans paths, so
// this helper exists only for naming simplicity.
func (c *Client) PathClean(p string) string {
	return strings.TrimSuffix(c.KeyClean(p), "/")
}

// KeyBase returns the last element of path k and cleans it. If k is empty or all slashes,
// return an empty string
func (c *Client) KeyBase(k string) string {
	addSlash := strings.HasSuffix(k, "/")

	base := path.Base(c.KeyJoin(k))
	if base == "/" || base == "." {
		return ""
	}
	if addSlash {
		return c.KeyJoin(base, "/")
	}
	return base
}

// PathBase calls KeyBase(p) and also trims the trailing '/' if necessary
func (c *Client) PathBase(p string) string {
	return strings.TrimSuffix(c.KeyBase(p), "/")
}

// SliceAddKeyPrefix takes in a slice of keys (strings) and a prefix and adds
// that prefix to every key in the slice
func (c *Client) SliceAddKeyPrefix(ss []string, p string) {
	for i, s := range ss {
		ss[i] = c.KeyJoin(p, s)
	}
}

// SliceTrimKeyPrefix takes in a slice of keys (strings) and a prefix and trims
// that prefix from every key in the slice
func (c *Client) SliceTrimKeyPrefix(ss []string, p string) {
	for i, s := range ss {
		ss[i] = c.KeyJoin(strings.TrimPrefix(s, p))
	}
}

// SliceRemoveFolders takes a list of keys and removes any folders (strings that end in a '/')
// and returns the new filtered list
func (c *Client) SliceRemoveFolders(ss []string) []string {
	var output []string
	for _, s := range ss {
		if !c.KeyIsFolder(s) {
			output = append(output, s)
		}
	}
	return output
}

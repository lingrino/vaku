package vaku

import (
	"path"
	"strings"
)

// KeyIsFolder returns true if the string ends in a '/'
func (c *Client) KeyIsFolder(k string) bool {
	if strings.HasSuffix(k, "/") {
		return true
	}
	return false
}

// KeyJoin takes n strings and combines them into a clean vault key
// vault keys have a trailing '/' if they are a folder
func (c *Client) KeyJoin(ks ...string) string {
	if strings.HasSuffix(ks[len(ks)-1], "/") {
		return strings.TrimPrefix(path.Join(ks...)+"/", "/")
	}
	return strings.TrimPrefix(path.Join(ks...), "/")
}

// PathJoin takes n strings and combines them into a clean vault path
// Paths in vault never begin or end in a slash
func (c *Client) PathJoin(ps ...string) string {
	return strings.TrimSuffix(c.KeyJoin(ps...), "/")
}

// KeyClean just calls KeyJoin with one string
// KeyJoin already cleans keys, this exists only for naming simplicity
func (c *Client) KeyClean(k string) string {
	return c.KeyJoin(k)
}

// PathClean just calls PathJoin with one string
// PathJoin already cleans paths, this exists only for naming simplicity
func (c *Client) PathClean(p string) string {
	return strings.TrimSuffix(c.KeyClean(p), "/")
}

// KeyBase returns the last element of k and cleans it
// If empty or all slashes, return an empty string
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

// PathBase returns the last element of p and cleans it
// If empty or all slashes, return an empty string
func (c *Client) PathBase(p string) string {
	return strings.TrimSuffix(c.KeyBase(p), "/")
}

// SliceAddKeyPrefix adds a prefix to every key in a slice
func (c *Client) SliceAddKeyPrefix(ss []string, p string) {
	for i, s := range ss {
		ss[i] = c.KeyJoin(p, s)
	}
}

// SliceTrimKeyPrefix trims a prefix from every key in a slice
func (c *Client) SliceTrimKeyPrefix(ss []string, p string) {
	for i, s := range ss {
		ss[i] = c.KeyJoin(strings.TrimPrefix(s, p))
	}
}

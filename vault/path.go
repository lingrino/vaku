package vault

import (
	pth "path"
	str "strings"
)

// PathIsFolder returns true if the string ends in a '/'
// Note that all 'cleaned' paths do not have a trailing slash
// and therefore are not folders.
func (c *Client) PathIsFolder(s string) bool {
	if str.HasSuffix(s, "/") {
		return true
	}
	return false
}

// PathJoin takes n strings and combines them into a clean vault path
func (c *Client) PathJoin(paths ...string) string {
	return str.TrimPrefix(pth.Join(paths...), "/")
}

// PathClean just calls PathJoin with one string
// PathJoin already cleans paths, this exists only for naming simplicity
func (c *Client) PathClean(path string) string {
	return c.PathJoin(path)
}

// PathLength returns the number of parts in a path
func (c *Client) PathLength(path string) int {
	path = c.PathJoin(path)
	if path == "" {
		return 0
	}
	return str.Count(path, "/") + 1
}

// PathGetPrefix returns the prefix of a path with depth n
// Ex: PathGetPrefix("secret/app/dev/env", 2) => "secret/app"
// If depth is negative or greater than the length or the path, does nothing"
func (c *Client) PathGetPrefix(path string, depth int) string {
	var pfx string

	pl := c.PathLength(path)
	if depth <= 0 || depth >= pl {
		return path
	}

	path = c.PathJoin(path)
	split := str.Split(path, "/")
	for i := 0; i < depth; i++ {
		pfx = c.PathJoin(pfx, split[i])
	}
	return pfx
}

// PathGetSuffix returns the suffix of a path with depth n
// Ex: PathGetSuffix("secret/app/dev/env", 2) => "dev/env"
// If depth is negative or greater than the length or the path, does nothing
func (c *Client) PathGetSuffix(path string, depth int) string {
	var sfx string

	pl := c.PathLength(path)
	if depth <= 0 || depth >= pl {
		return path
	}

	path = c.PathJoin(path)
	for i := 0; i < depth; i++ {
		sfx = c.PathJoin(pth.Base(path), sfx)
		path = pth.Dir(path)
	}
	return sfx
}

// PathRemovePrefix takes a path and returns it with n leading
// parts removed, where n is the depth specified
// If depth is negative or greater than the length or the path, does nothing
func (c *Client) PathRemovePrefix(path string, depth int) string {
	pl := c.PathLength(path)
	if depth <= 0 || depth >= pl {
		return path
	}

	return c.PathGetSuffix(path, pl-depth)
}

// PathRemoveSuffix takes a path and returns it with n trailing
// parts removed, where n is the depth specified
// If depth is negative or greater than the length or the path, does nothing
func (c *Client) PathRemoveSuffix(path string, depth int) string {
	pl := c.PathLength(path)
	if depth <= 0 || depth >= pl {
		return path
	}

	return c.PathGetPrefix(path, pl-depth)
}

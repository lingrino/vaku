package vaku2

import (
	"path"
	"strings"
)

// KeyJoin combines strings into a clean Vault key. Keys may have a trailing '/' to signify they are a folder.
func KeyJoin(k ...string) string {
	if strings.HasSuffix(k[len(k)-1], "/") {
		return strings.TrimPrefix(path.Join(k...)+"/", "/")
	}
	return strings.TrimPrefix(path.Join(k...), "/")
}

// PathJoin combines strings into a clean Vault path. Paths never end with a '/'.
func PathJoin(p ...string) string {
	return strings.TrimSuffix(KeyJoin(p...), "/")
}

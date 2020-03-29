package vaku

import (
	"path"
	"strings"
)

// IsFolder if path is a folder (ends in "/").
func IsFolder(p string) bool {
	return strings.HasSuffix(p, "/")
}

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

// PrefixList adds a prefix to every item in a list
func PrefixList(list []string, prefix string) {
	for i, v := range list {
		list[i] = KeyJoin(prefix, v)
	}
}

// TrimListPrefix removes a prefix from every item in a list
func TrimListPrefix(list []string, prefix string) {
	for i, v := range list {
		list[i] = KeyJoin(strings.TrimPrefix(v, prefix))
	}
}

package vaku

import (
	"path"
	"strings"
)

// PathJoin combines multiple paths into one.
func PathJoin(p ...string) string {
	if strings.HasSuffix(p[len(p)-1], "/") {
		return strings.TrimPrefix(path.Join(p...)+"/", "/")
	}
	return strings.TrimPrefix(path.Join(p...), "/")
}

// IsFolder if path is a folder (ends in "/").
func IsFolder(p string) bool {
	return strings.HasSuffix(p, "/")
}

// EnsureFolder ensures a path is a folder (adds a trailing "/").
func EnsureFolder(p string) string {
	return PathJoin(p, "/")
}

// AddPrefix adds a prefix to a path.
func AddPrefix(p, pfx string) string {
	return PathJoin(pfx, p)
}

// EnsurePrefix adds a prefix to a path if it doesn't already have it.
func EnsurePrefix(p, pfx string) string {
	if strings.HasPrefix(p, pfx) {
		return p
	}
	return PathJoin(pfx, p)
}

// AddPrefixList adds a prefix to every item in a list.
func AddPrefixList(l []string, pfx string) {
	for i, v := range l {
		l[i] = PathJoin(pfx, v)
	}
}

// EnsurePrefixList adds a prefix to every item in a list if it doesn't already have it.
func EnsurePrefixList(l []string, pfx string) {
	for i, v := range l {
		if !strings.HasPrefix(v, pfx) {
			l[i] = PathJoin(pfx, v)
		}
	}
}

// TrimPrefixList removes a prefix from every item in a list.
func TrimPrefixList(l []string, pfx string) {
	for i, v := range l {
		l[i] = PathJoin(strings.TrimPrefix(v, pfx))
	}
}

// EnsurePrefixMap ensures a prefix for every key in a map.
func EnsurePrefixMap(m map[string]map[string]interface{}, pfx string) {
	for k, v := range m {
		delete(m, k)
		m[EnsurePrefix(k, pfx)] = v
	}
}

// TrimPrefixMap removes a prefix from every key in a map.
func TrimPrefixMap(m map[string]map[string]interface{}, pfx string) {
	for k, v := range m {
		delete(m, k)
		m[PathJoin(strings.TrimPrefix(k, pfx))] = v
	}
}

// InsertIntoPath adds 'insert' into 'path' after 'after' and returns the new path.
func InsertIntoPath(path, after, insert string) string {
	return PathJoin(after, insert, strings.TrimPrefix(path, after))
}

// errFuncOnChan takes a function like errgroup.Wait() and provides a channel that can be read for
// the err value that the function returns. Makes it easy to wait inside of a select statement.
func errFuncOnChan(errFunc func() error) <-chan error {
	errC := make(chan error)
	go func() {
		errC <- errFunc()
		close(errC)
	}()
	return errC
}

// mergeMaps merges m2 into m1, preferring data from m2.
func mergeMaps(m1, m2 map[string]map[string]interface{}) {
	for k, v := range m2 {
		m1[k] = v
	}
}

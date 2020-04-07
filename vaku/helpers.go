package vaku

import (
	"path"
	"strings"
)

// IsFolder if path is a folder (ends in "/").
func IsFolder(p string) bool {
	return strings.HasSuffix(p, "/")
}

// MakeFolder adds a slash to the end of a path if it doesn't already have one
func MakeFolder(p string) string {
	return KeyJoin(p, "/")
}

// EnsurePrefix adds a prefix to a path if it doesn't already have it
func EnsurePrefix(p, pfx string) string {
	if strings.HasPrefix(p, pfx) {
		return p
	}
	return KeyJoin(pfx, p)
}

// KeyJoin combines strings into a clean Vault key. Keys may have a trailing '/' to signify they are
// a folder.
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

// TrimMapKeyPrefix removes a prefix from every key in a map
func TrimMapKeyPrefix(m map[string]map[string]interface{}, prefix string) {
	for k, v := range m {
		delete(m, k)
		m[KeyJoin(strings.TrimPrefix(k, prefix))] = v
	}
}

// waitFuncOnChan takes a function like waitgroup.Wait() and provides a channel that can be read
// after the function returns. Makes it easy to wait inside of a select statement.
func waitFuncOnChan(waitFunc func()) <-chan bool {
	done := make(chan bool)
	go func() {
		waitFunc()
		done <- true
		close(done)
	}()
	return done
}

// errFuncOnChan takes a function like errgroup.Wait() and provides a channel that can be read for the err value
// that the function returns. Makes it easy to wait inside of a select statement.
func errFuncOnChan(errFunc func() error) <-chan error {
	errC := make(chan error)
	go func() {
		errC <- errFunc()
		close(errC)
	}()
	return errC
}

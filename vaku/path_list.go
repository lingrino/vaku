package vaku

import (
	"errors"
	"sort"
)

var (
	ErrPathList  = errors.New("path list")
	ErrVaultList = errors.New("vault list")
)

// PathList lists paths at a path.
func (c *Client) PathList(p string) ([]string, error) {
	return c.pathList(c.srcL, p)
}

// PathListDest lists paths at a path.
func (c *Client) PathListDst(p string) ([]string, error) {
	return c.pathList(c.dstL, p)
}

// pathList does the actual list.
func (c *Client) pathList(l logical, p string) ([]string, error) {
	secret, err := l.List(p)
	if err != nil {
		return nil, newWrapErr(p, ErrPathList, newWrapErr(err.Error(), ErrVaultList, nil))
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	data, ok := secret.Data["keys"]
	if !ok || data == nil {
		return nil, newWrapErr(p, ErrPathList, ErrDecodeSecret)
	}
	keys, ok := data.([]interface{})
	if !ok {
		return nil, newWrapErr(p, ErrPathList, ErrDecodeSecret)
	}

	output := make([]string, len(keys))
	for i, k := range keys {
		key, ok := k.(string)
		if !ok {
			return nil, newWrapErr(p, ErrPathList, ErrDecodeSecret)
		}
		output[i] = key
	}

	if c.fullPath {
		PrefixList(output, p)
	}

	sort.Strings(output)
	return output, nil
}

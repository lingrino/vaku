package vaku

import (
	"errors"
	"sort"

	"github.com/hashicorp/vault/api"
)

var (
	// ErrPathList when PathList/PathListDst fails.
	ErrPathList = errors.New("path list")
	// ErrVaultList when the underlying Vault API list fails.
	ErrVaultList = errors.New("vault list")
)

// PathList lists paths at a path.
func (c *Client) PathList(p string) ([]string, error) {
	return c.pathList(c.srcL, p)
}

// PathListDst lists paths at a path.
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

	list, err := decodeSecret(secret)
	if err != nil {
		return nil, newWrapErr(p, ErrPathList, err)
	}

	if c.absolutepath {
		PrefixList(list, p)
	}

	sort.Strings(list)
	return list, nil
}

func decodeSecret(secret *api.Secret) ([]string, error) {
	data, ok := secret.Data["keys"]
	if !ok || data == nil {
		return nil, newWrapErr("", ErrDecodeSecret, nil)
	}
	keys, ok := data.([]interface{})
	if !ok {
		return nil, newWrapErr("", ErrDecodeSecret, nil)
	}

	output := make([]string, len(keys))
	for i, k := range keys {
		key, ok := k.(string)
		if !ok {
			return nil, newWrapErr("", ErrDecodeSecret, nil)
		}
		output[i] = key
	}

	return output, nil
}

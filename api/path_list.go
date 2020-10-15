package vaku

import (
	"errors"

	vault "github.com/hashicorp/vault/api"
)

var (
	// ErrPathList when PathList fails.
	ErrPathList = errors.New("path list")
	// ErrVaultList when the underlying Vault API list fails.
	ErrVaultList = errors.New("vault list")
)

// PathList lists paths at a path.
func (c *Client) PathList(p string) ([]string, error) {
	vaultPath, _, err := c.rewritePath(p, vaultList)
	if err != nil {
		return nil, newWrapErr(p, ErrPathList, err)
	}

	secret, err := c.vl.List(vaultPath)
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

	c.outputPaths(list, p)

	return list, nil
}

func decodeSecret(secret *vault.Secret) ([]string, error) {
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

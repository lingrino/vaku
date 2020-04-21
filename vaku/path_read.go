package vaku

import (
	"errors"
)

var (
	// ErrPathRead when PathRead fauls.
	ErrPathRead = errors.New("path read")
	// ErrVaultRead when the underlying Vault API read fails.
	ErrVaultRead = errors.New("vault read")
)

// PathRead gets data at a path.
func (c *Client) PathRead(p string) (map[string]interface{}, error) {
	secret, err := c.vl.Read(p)
	if err != nil {
		return nil, newWrapErr(p, ErrPathRead, newWrapErr(err.Error(), ErrVaultRead, nil))
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	return secret.Data, nil
}

package vaku

import (
	"errors"
)

var (
	// ErrPathWrite when PathWrite errors.
	ErrPathWrite = errors.New("path write")
	// ErrVaultWrite when the underlying Vault API write fails.
	ErrVaultWrite = errors.New("vault write")
)

// PathWrite writes data to a path.
func (c *Client) PathWrite(p string, d map[string]interface{}) error {
	if d == nil {
		return newWrapErr(p, ErrPathWrite, ErrNilData)
	}

	_, err := c.vl.Write(p, d)
	if err != nil {
		return newWrapErr(p, ErrPathWrite, newWrapErr(err.Error(), ErrVaultWrite, nil))
	}

	return nil
}

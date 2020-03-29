package vaku

import (
	"errors"
)

var (
	// ErrPathWrite when PathWrite/PathWriteDest errors
	ErrPathWrite = errors.New("path write")
	// ErrVaultWrite when the underlying Vault API write fails
	ErrVaultWrite = errors.New("vault write")
)

// PathWrite writes data to a path.
func (c *Client) PathWrite(p string, d map[string]interface{}) error {
	return c.pathWrite(c.srcL, p, d)
}

// PathWriteDst writes data to a path.
func (c *Client) PathWriteDst(p string, d map[string]interface{}) error {
	return c.pathWrite(c.dstL, p, d)
}

// pathWrite does the actual write.
func (c *Client) pathWrite(l logical, p string, d map[string]interface{}) error {
	if d == nil {
		return newWrapErr(p, ErrPathWrite, ErrNilData)
	}

	_, err := l.Write(p, d)
	if err != nil {
		return newWrapErr(p, ErrPathWrite, newWrapErr(err.Error(), ErrVaultWrite, nil))
	}

	return nil
}

package vaku2

import (
	"errors"
)

var (
	ErrPathRead  = errors.New("path read")
	ErrVaultRead = errors.New("vault read")
)

// PathRead gets data at a path.
func (c *Client) PathRead(p string) (map[string]interface{}, error) {
	return c.pathRead(c.srcL, p)
}

// PathReadDest gets data at a path.
func (c *Client) PathReadDst(p string) (map[string]interface{}, error) {
	return c.pathRead(c.dstL, p)
}

// pathRead does the actual read.
func (c *Client) pathRead(l logical, p string) (map[string]interface{}, error) {
	secret, err := l.Read(p)
	if err != nil {
		return nil, newWrapErr(p, ErrPathRead, newWrapErr(err.Error(), ErrVaultRead, nil))
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	return secret.Data, nil
}

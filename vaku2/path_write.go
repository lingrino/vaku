package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrVaultWrite = errors.New("vault write")
)

// PathWrite writes data to a path using the source client.
func (c *Client) PathWrite(p string, d map[string]interface{}) error {
	return c.pathWrite(c.sourceL, p, d)
}

// PathWriteDest writes data to a path using the dest client.
func (c *Client) PathWriteDest(p string, d map[string]interface{}) error {
	return c.pathWrite(c.destL, p, d)
}

// pathWrite writes data to a path.
func (c *Client) pathWrite(apiL logical, p string, d map[string]interface{}) error {
	_, err := apiL.Write(p, d)
	if err != nil {
		return newWrapErr(fmt.Sprintf("%q: %v: %v", p, ErrVaultWrite, err), ErrVaultWrite, nil)
	}

	return nil
}

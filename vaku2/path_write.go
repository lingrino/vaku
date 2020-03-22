package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrVaultWrite = errors.New("vault write")
)

// PathWrite writes data to a path.
func (c *Client) PathWrite(p string, d map[string]interface{}) error {
	_, err := c.sourceL.Write(p, d)
	if err != nil {
		return fmt.Errorf("%q: %w", p, ErrVaultWrite)
	}

	return nil
}

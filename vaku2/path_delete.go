package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrVaultDelete = errors.New("vault delete")
)

// PathDelete deletes data at a path
func (c *Client) PathDelete(p string) error {
	_, err := c.sourceL.Delete(p)
	if err != nil {
		return fmt.Errorf("%q: %w", p, ErrVaultDelete)
	}

	return nil
}

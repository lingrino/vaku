package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrVaultRead = errors.New("vault list")
)

// PathRead takes a path, calls vault read, extracts the secret, and returns it.
func (c *Client) PathRead(p string) (map[string]interface{}, error) {
	secret, err := c.sourceL.Read(p)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", p, ErrVaultRead)
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	return secret.Data, nil
}

package vaku2

import (
	"errors"
	"fmt"
	"sort"
)

var (
	ErrVaultList    = errors.New("vault list")
	ErrDecodeSecret = errors.New("decode secret")
)

// PathList takes a path, calls vault list, extracts the secret as a list of keys, and returns it.
func (c *Client) PathList(p string) ([]string, error) {
	secret, err := c.sourceL.List(p)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", p, ErrVaultList)
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	data, ok := secret.Data["keys"]
	if !ok || data == nil {
		return nil, fmt.Errorf("%w", ErrDecodeSecret)
	}
	keys, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w", ErrDecodeSecret)
	}

	output := make([]string, len(keys))
	for i, k := range keys {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("%w", ErrDecodeSecret)
		}
		output[i] = key
	}

	if c.fullPath {
		PrefixList(output, p)
	}

	sort.Strings(output)
	return output, nil
}

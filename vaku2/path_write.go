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
	fmt.Printf("%+v\n", d)
	_, err := apiL.Write(p, d)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%q: %w", p, ErrVaultWrite)
	}

	return nil
}

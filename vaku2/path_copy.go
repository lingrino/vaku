package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrPathCopy = errors.New("path copy")
)

// PathCopy reads a secret at a source path and copies it to the destination path. When copying from
// one vault server to another the client must have been initialized using WithDestClient().
func (c *Client) PathCopy(source, dest string) error {
	secret, err := c.PathRead(source)
	if err != nil {
		return newWrapErr(fmt.Sprintf("%v: %v", ErrPathCopy, err), ErrPathCopy, err)
	}

	err = c.PathWriteDest(dest, secret)
	if err != nil {
		return newWrapErr(fmt.Sprintf("%v: %v", ErrPathCopy, err), ErrPathCopy, err)
	}

	return nil
}

package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrPathMove = errors.New("path move")
)

// PathMove calls PathCopy and then deletes the source path. When copying from one vault server to
// another the client must have been initialized using WithDestClient().
func (c *Client) PathMove(source, dest string) error {
	err := c.PathCopy(source, dest)
	if err != nil {
		return newWrapErr(fmt.Sprintf("%v: %v", ErrPathMove, err), ErrPathMove, err)
	}

	err = c.PathDelete(source)
	if err != nil {
		return newWrapErr(fmt.Sprintf("%v: %v", ErrPathMove, err), ErrPathMove, err)
	}

	return nil
}

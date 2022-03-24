package vaku

import (
	"errors"
)

var (
	// ErrPathUpdate when PathUpdate fails.
	ErrPathUpdate = errors.New("path update")
)

// PathUpdate updates a path with data. New data (precedence) is merged with existing data.
func (c *Client) PathUpdate(p string, d map[string]any) error {
	if d == nil {
		return newWrapErr(p, ErrPathUpdate, ErrNilData)
	}

	read, err := c.PathRead(p)
	if err != nil {
		return newWrapErr(p, ErrPathUpdate, err)
	}
	if read == nil {
		read = make(map[string]any, len(d))
	}

	for k, v := range d {
		read[k] = v
	}

	err = c.PathWrite(p, read)
	if err != nil {
		return newWrapErr(p, ErrPathUpdate, err)
	}

	return nil
}

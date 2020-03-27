package vaku2

import (
	"errors"
)

var (
	ErrPathUpdate = errors.New("path update")
)

// PathUpdate updates a path with data. Existing data is merged with new data. Precedence is given
// to new data.
func (c *Client) PathUpdate(p string, d map[string]interface{}) error {
	if d == nil {
		return newWrapErr(p, ErrPathUpdate, ErrNilData)
	}

	read, err := c.PathRead(p)
	if err != nil {
		return newWrapErr(p, ErrPathUpdate, err)
	}
	if read == nil {
		read = make(map[string]interface{}, len(d))
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

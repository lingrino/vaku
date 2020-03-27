package vaku2

import (
	"errors"
	"fmt"
)

var (
	ErrPathUpdate = errors.New("path update")
)

// PathUpdate takes a path with existing data and new data to write to that path. It merges the data
// at the existing path with the new data and writes the merged data back to Vault. Precedence is
// given to the new data.
func (c *Client) PathUpdate(p string, d map[string]interface{}) error {
	if d == nil {
		return newWrapErr(fmt.Sprintf("%v", ErrPathUpdate), ErrPathUpdate, ErrNilData)
	}

	read, err := c.PathRead(p)
	if err != nil {
		return newWrapErr(fmt.Sprintf("read %q: %v: %v", p, ErrPathUpdate, err), ErrPathUpdate, err)
	}
	if read == nil {
		read = make(map[string]interface{}, len(d))
	}

	for k, v := range d {
		read[k] = v
	}

	err = c.PathWrite(p, read)
	if err != nil {
		return newWrapErr(fmt.Sprintf("write %q: %v: %v", p, ErrPathUpdate, err), ErrPathUpdate, err)
	}

	return nil
}

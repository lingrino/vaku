package vaku

import (
	"github.com/pkg/errors"
)

// PathSearch takes in a path and search string, reads the path,
// and searches the data for a match. Returns true if the data is found
func (c *Client) PathSearch(i *PathInput, s string) (bool, error) {
	var err error

	// Read the data at the path
	read, err := c.PathRead(i)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to read data at path %s", i.Path)
	}

	return false, err
}

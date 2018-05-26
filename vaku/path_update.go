package vaku

import (
	"github.com/pkg/errors"
)

// PathUpdate takes in a path and new data to write. If then dates the data
// at the existing path with the new data, merging the two with precedence given
// to the new data.
func (c *Client) PathUpdate(i *PathInput, d map[string]interface{}) error {
	var err error

	// Get old data
	read, err := c.PathRead(i)
	if err != nil {
		return errors.Wrapf(err, "Failed to read data at path %s. PathUpdate only works on existing data", i.Path)
	}

	// Generate the new data to write
	for k, v := range d {
		read[k] = v
	}

	// Write the updated data back to vault
	err = c.PathWrite(i, read)
	if err != nil {
		return errors.Wrapf(err, "Failed to write updated data back to %s", i.opPath)
	}

	return err
}

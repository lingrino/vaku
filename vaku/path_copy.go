package vaku

import (
	"github.com/pkg/errors"
)

// PathCopy takes in a source PathInput and a target PathInput.
// It then copies the data from one path to another. PathCopy can
// copy from one mount to another by default. Note that this will
// overwrite an existing key at the target path.
func (c *Client) PathCopy(s *PathInput, t *PathInput) error {
	var err error

	// Read the data from the source path
	d, err := c.PathRead(s)
	if err != nil {
		return errors.Wrapf(err, "Failed to read data at %s", s.Path)
	}

	// Write the data to the new path
	err = c.PathWrite(t, d)
	if err != nil {
		return errors.Wrapf(err, "Failed to write data to %s", t.Path)
	}

	return err
}

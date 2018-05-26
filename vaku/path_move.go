package vaku

import (
	"github.com/pkg/errors"
)

// PathMove takes in a source PathInput and a target PathInput.
// It then moves (destructively) the data from one path to another.
// PathMove can move from one mount to another by default. Note that
// this will overwrite any existing key at the target path.
func (c *Client) PathMove(s *PathInput, t *PathInput) error {
	var err error

	// Copy the data to the new path
	err = c.PathCopy(s, t)
	if err != nil {
		return errors.Wrapf(err, "Failed to copy data from %s to %s", s.Path, t.Path)
	}

	// Delete the data at the old path
	err = c.PathDelete(s)
	if err != nil {
		return errors.Wrapf(err, "Failed to delete source path %s. This means that the path was copied instead of deleted", s.Path)
	}

	return err
}

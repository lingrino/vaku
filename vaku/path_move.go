package vaku

import (
	"github.com/pkg/errors"
)

// PathMove calls PathCopy() with the same inputs followed by PathDelete() on
// the source if the copy was successful. Note that this will overwrite any existing
// keys at the target Path.
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

package vaku

import (
	"github.com/pkg/errors"
)

// PathDelete takes in a path and deletes it. For v2
// mounts this function only "marks the path as deleted",
// it does nothing with the versions of the path
func (c *Client) PathDelete(i *PathInput) error {
	var err error

	// Initialize the input
	i.opType = "delete"
	err = c.InitPathInput(i)
	if err != nil {
		return errors.Wrapf(err, "Failed to init delete path %s", i.Path)
	}

	// Do the actual delete
	_, err = c.Logical().Delete(i.opPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to delete secret at %s", i.opPath)
	}

	return err
}

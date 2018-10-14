package vaku

import (
	"github.com/pkg/errors"
)

// PathDestroy takes in a PathInput and calls the native delete on 'mount/metadata/path'
// This function only works on versioned (V2) key/value mounts. Note that this destroys ALL
// versions at the path, there is no current support for destroying specific versions.
func (c *Client) PathDestroy(i *PathInput) error {
	var err error

	// Initialize the input
	i.opType = "destroy"
	err = c.InitPathInput(i)
	if err != nil {
		return errors.Wrapf(err, "Failed to init destroy path %s", i.Path)
	}

	// Do the actual destroy
	_, err = c.Logical().Delete(i.opPath)
	if err != nil {
		return errors.Wrapf(err, "Failed to destroy secret at %s", i.opPath)
	}

	return err
}

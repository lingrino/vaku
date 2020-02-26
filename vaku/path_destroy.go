package vaku

import "fmt"

// PathDestroy takes in a PathInput and calls the native delete on 'mount/metadata/path'
// This function only works on versioned (V2) key/value mounts. Note that this destroys ALL
// versions at the path, there is no current support for destroying specific versions.
func (c *Client) PathDestroy(i *PathInput) error {
	var err error

	// Initialize the input
	i.opType = "destroy"
	err = c.InitPathInput(i)
	if err != nil {
		return fmt.Errorf("failed to init destroy path %s: %w", i.Path, err)
	}

	// Do the actual destroy
	_, err = c.Logical().Delete(i.opPath)
	if err != nil {
		return fmt.Errorf("failed to destroy secret at %s: %w", i.opPath, err)
	}

	return err
}

package vaku

import "fmt"

// PathDestroyVersions takes in a PathInput and calls the destroy on 'mount/destroy/path'
// This function only works on versioned (V2) key/value mounts.
func (c *Client) PathDestroyVersions(i *PathInput, versions []int) error {
	var err error

	// Initialize the input
	i.opType = "destroyversions"
	err = c.InitPathInput(i)
	if err != nil {
		return fmt.Errorf("failed to init destroy path %s: %w", i.Path, err)
	}

	d := map[string]interface{}{
		"versions": versions,
	}

	// Do the actual destroy
	_, err = c.Logical().Write(i.opPath, d)
	if err != nil {
		return fmt.Errorf("failed to destroy secret at %s: %w", i.opPath, err)
	}

	return err
}

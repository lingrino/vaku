package vaku

import "fmt"

// PathDelete takes in a PathInput and calls the native Vault delete on it. For v2
// mounts this function only "marks the path as deleted", it does nothing with the
// versions of the path
func (c *Client) PathDelete(i *PathInput) error {
	var err error

	// Initialize the input
	i.opType = "delete"
	err = c.InitPathInput(i)
	if err != nil {
		return fmt.Errorf("failed to init delete path %s: %w", i.Path, err)
	}

	// Do the actual delete
	_, err = c.Logical().Delete(i.opPath)
	if err != nil {
		return fmt.Errorf("failed to delete secret at %s: %w", i.opPath, err)
	}

	return err
}

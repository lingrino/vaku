package vaku

import "fmt"

// PathUpdate takes in a path with existing data and new data to write to that path.
// It then merges the data at the existing path with the new data, with precedence given
// to the new data, and writes the merged data back to Vault
func (c *Client) PathUpdate(i *PathInput, d map[string]interface{}) error {
	var err error

	// Get old data
	read, err := c.PathRead(i)
	if err != nil {
		return fmt.Errorf("failed to read data at path %s. PathUpdate only works on existing data: %w", i.Path, err)
	}

	// Generate the new data to write
	for k, v := range d {
		read[k] = v
	}

	// Write the updated data back to vault
	err = c.PathWrite(i, read)
	if err != nil {
		return fmt.Errorf("failed to write updated data back to %s: %w", i.opPath, err)
	}

	return err
}

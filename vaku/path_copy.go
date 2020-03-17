package vaku

import "fmt"

// PathCopy takes in a source PathInput and a target PathInput. It then copies the data
// from one path to another. Note that PathCopy can be used to copy data from one mount
// to another. Note also that this will overwrite any existing key at the target path.
func (c *Client) PathCopy(s *PathInput, t *PathInput) error {
	var err error

	// Read the data from the source path
	d, err := c.PathRead(s)
	if err != nil {
		return fmt.Errorf("failed to read data at %s: %w", s.Path, err)
	}

	// Write the data to the new path
	err = c.PathWrite(t, d)
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %w", t.Path, err)
	}

	return err
}

// PathCopy takes in a source PathInput and a target PathInput. It then copies the data
// from one path to another. Note that PathCopy can be used to copy data from one mount
// to another. Note also that this will overwrite any existing key at the target path.
func (c *CopyClient) PathCopy(s *PathInput, t *PathInput) error {
	var err error

	// Read the data from the source path
	d, err := c.Source.PathRead(s)
	if err != nil {
		return fmt.Errorf("failed to read data at %s: %w", s.Path, err)
	}

	// Do not copy KV v2 secrets that are deleted
	if s.mountVersion == "2" && (d["VAKU_STATUS"] == "SECRET_HAS_BEEN_DELETED" ||
		d["VAKU_STATUS"] == "SECRET_HAS_BEEN_DESTROYED") {
		return nil
	}

	// Write the data to the new path
	err = c.Target.PathWrite(t, d)
	if err != nil {
		return fmt.Errorf("failed to write data to %s: %w", t.Path, err)
	}

	return err
}

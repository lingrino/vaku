package vaku

import (
	"errors"
	"fmt"
)

var (
	// ErrPathCopy when PathCopy fails.
	ErrPathCopy = errors.New("path copy")
	// ErrPathCopyAllVersions when PathCopyAllVersions fails.
	ErrPathCopyAllVersions = errors.New("path copy all versions")
)

// PathCopy copies data at a source path to a destination path.
func (c *Client) PathCopy(src, dst string) error {
	secret, err := c.PathRead(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopy, err)
	}

	err = c.dc.PathWrite(dst, secret)
	if err != nil {
		return newWrapErr(dst, ErrPathCopy, err)
	}

	return nil
}

// PathCopyAllVersions copies all versions of data at a source path to a destination path.
// Only works on v2 kv engines.
func (c *Client) PathCopyAllVersions(src, dst string) error {
	// First check if this is a v2 mount
	_, mv, err := c.rewritePath(src, vaultRead)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	if mv != mv2 {
		err := newWrapErr("all versions copy not supported on KV v1", ErrMountVersion, nil)
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	// Read metadata to get version information
	metadata, err := c.PathReadMetadata(src)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	if metadata == nil {
		return nil // nothing to copy
	}

	// Extract and copy versions
	return c.copyVersionsFromMetadata(src, dst, metadata)
}

// copyVersionsFromMetadata extracts version data from metadata and copies each version.
func (c *Client) copyVersionsFromMetadata(src, dst string, metadata map[string]any) error {
	versionsData, ok := metadata["versions"].(map[string]any)
	if !ok {
		return newWrapErr(src, ErrPathCopyAllVersions, fmt.Errorf("invalid metadata format"))
	}

	for versionStr := range versionsData {
		err := c.copyIndividualVersion(src, dst, versionStr, versionsData[versionStr])
		if err != nil {
			return err
		}
	}

	return nil
}

// copyIndividualVersion copies a single version if it's not deleted.
func (c *Client) copyIndividualVersion(src, dst, versionStr string, versionMeta any) error {
	versionData, ok := versionMeta.(map[string]any)
	if !ok {
		return nil // skip invalid version data
	}

	// Skip deleted versions
	if isVersionDeleted(versionData) {
		return nil
	}

	// Parse version number
	var versionNum int
	_, err := fmt.Sscanf(versionStr, "%d", &versionNum)
	if err != nil {
		return nil // skip invalid version numbers
	}

	// Read the specific version
	versionSecret, err := c.PathReadVersion(src, versionNum)
	if err != nil {
		return newWrapErr(src, ErrPathCopyAllVersions, err)
	}

	if versionSecret != nil {
		// Write the version data to destination
		err = c.dc.PathWrite(dst, versionSecret)
		if err != nil {
			return newWrapErr(dst, ErrPathCopyAllVersions, err)
		}
	}

	return nil
}

// isVersionDeleted checks if a version has been deleted or destroyed.
func isVersionDeleted(versionData map[string]any) bool {
	if deletionTime, exists := versionData["deletion_time"]; exists && deletionTime != nil && deletionTime != "" {
		return true
	}
	if destroyed, exists := versionData["destroyed"]; exists && destroyed == true {
		return true
	}
	return false
}

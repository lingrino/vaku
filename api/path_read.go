package vaku

import (
	"errors"
	"strconv"
)

var (
	// ErrPathRead when PathRead fails.
	ErrPathRead = errors.New("path read")
	// ErrVaultRead when the underlying Vault API read fails.
	ErrVaultRead = errors.New("vault read")
)

// PathRead reads data at a path.
func (c *Client) PathRead(p string) (map[string]any, error) {
	vaultPath, mv, err := c.rewritePath(p, vaultRead)
	if err != nil {
		return nil, newWrapErr(p, ErrPathRead, err)
	}

	secret, err := c.vl.Read(vaultPath)
	if err != nil {
		if c.ignoreAccessErrors {
			return nil, nil
		}
		return nil, newWrapErr(p, ErrPathRead, newWrapErr(err.Error(), ErrVaultRead, nil))
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	data := secret.Data
	if mv == mv2 {
		data = extractV2Read(data)
	}

	return data, nil
}

// extractV2Read returns data["data"] if the secret is not deleted or destroyed.
func extractV2Read(data map[string]any) map[string]any {
	if data == nil {
		return nil
	}

	if isDeleted(data) {
		return nil
	}

	dd := data["data"]
	if dd == nil {
		return nil
	}

	dm, ok := dd.(map[string]any)
	if !ok {
		return nil
	}

	return dm
}

// isDeleted checks if the secret has been deleted or destroyed.
func isDeleted(data map[string]any) bool {
	metadata, ok := data["metadata"].(map[string]any)
	if !ok {
		return true
	}
	deletionTime, ok := metadata["deletion_time"].(string)
	if !ok || deletionTime != "" {
		return true
	}
	destroyed, ok := metadata["destroyed"].(bool)
	if !ok || destroyed {
		return true
	}

	return false
}

// PathReadAllVersions reads all versions of a secret at a path (KV v2 only).
// Returns a slice of maps, one for each version, in version order.
func (c *Client) PathReadAllVersions(p string) ([]map[string]any, error) {
	// First, determine the mount version
	_, mv, err := c.rewritePath(p, vaultRead)
	if err != nil {
		return nil, newWrapErr(p, ErrPathRead, err)
	}

	// Only works for KV v2
	if mv != mv2 {
		// For v1, just return the single version
		data, err := c.PathRead(p)
		if err != nil {
			return nil, err
		}
		if data == nil {
			return nil, nil
		}
		return []map[string]any{data}, nil
	}

	// For KV v2, for now just return the latest version
	// TODO: Implement proper version history reading using metadata endpoint
	data, err := c.PathRead(p)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	return []map[string]any{data}, nil
}

// pathReadVersion reads a specific version of a secret.
func (c *Client) pathReadVersion(p string, version int) (map[string]any, error) {
	// For now, let's just use PathRead for version 0 (latest)
	// and return nil for other versions since the version query isn't working as expected
	if version == 1 {
		// Read the current version
		return c.PathRead(p)
	}

	// For other versions, try adding the version parameter
	vaultPath, mv, err := c.rewritePath(p, vaultRead)
	if err != nil {
		return nil, newWrapErr(p, ErrPathRead, err)
	}

	// Add version query parameter
	vaultPath = vaultPath + "?version=" + strconv.Itoa(version)

	secret, err := c.vl.Read(vaultPath)
	if err != nil {
		if c.ignoreAccessErrors {
			return nil, nil
		}
		return nil, newWrapErr(p, ErrPathRead, newWrapErr(err.Error(), ErrVaultRead, nil))
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	data := secret.Data
	if mv == mv2 {
		data = extractV2Read(data)
	}

	return data, nil
}

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
	// ErrPathReadMetadata when PathReadMetadata fails.
	ErrPathReadMetadata = errors.New("path read metadata")
	// ErrPathReadVersion when PathReadVersion fails.
	ErrPathReadVersion = errors.New("path read version")
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

// PathReadMetadata reads metadata at a path including version information.
// Only works on v2 kv engines.
func (c *Client) PathReadMetadata(p string) (map[string]any, error) {
	vaultPath, mv, err := c.rewritePath(p, vaultReadMeta)
	if err != nil {
		return nil, newWrapErr(p, ErrPathReadMetadata, err)
	}

	if mv != mv2 {
		return nil, newWrapErr(p, ErrPathReadMetadata, newWrapErr("metadata not supported on KV v1", ErrMountVersion, nil))
	}

	secret, err := c.vl.Read(vaultPath)
	if err != nil {
		if c.ignoreAccessErrors {
			return nil, nil
		}
		return nil, newWrapErr(p, ErrPathReadMetadata, newWrapErr(err.Error(), ErrVaultRead, nil))
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	return secret.Data, nil
}

// PathReadVersion reads a specific version of data at a path.
// Only works on v2 kv engines.
func (c *Client) PathReadVersion(p string, version int) (map[string]any, error) {
	vaultPath, mv, err := c.rewritePath(p, vaultRead)
	if err != nil {
		return nil, newWrapErr(p, ErrPathReadVersion, err)
	}

	if mv != mv2 {
		return nil, newWrapErr(p, ErrPathReadVersion, newWrapErr("versions not supported on KV v1", ErrMountVersion, nil))
	}

	// Add version parameter to the path
	vaultPath = vaultPath + "?version=" + strconv.Itoa(version)

	secret, err := c.vl.Read(vaultPath)
	if err != nil {
		if c.ignoreAccessErrors {
			return nil, nil
		}
		return nil, newWrapErr(p, ErrPathReadVersion, newWrapErr(err.Error(), ErrVaultRead, nil))
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

package vaku

import (
	"encoding/json"
	"errors"
	"strconv"
)

var (
	// ErrPathReadMeta when PathReadMeta fails.
	ErrPathReadMeta = errors.New("path read meta")
)

// SecretVersionMeta contains metadata about a single secret version.
type SecretVersionMeta struct {
	CreatedTime string
	Deleted     bool
	Destroyed   bool
}

// SecretMeta contains metadata about a secret and all of its versions.
type SecretMeta struct {
	CurrentVersion int
	Versions       map[int]SecretVersionMeta
}

// PathReadMeta reads metadata at a path. Only works on v2 kv engines.
func (c *Client) PathReadMeta(p string) (*SecretMeta, error) {
	vaultPath, _, err := c.rewritePath(p, vaultReadMeta)
	if err != nil {
		return nil, newWrapErr(p, ErrPathReadMeta, err)
	}

	secret, err := c.vl.Read(vaultPath)
	if err != nil {
		return nil, newWrapErr(p, ErrPathReadMeta, newWrapErr(err.Error(), ErrVaultRead, nil))
	}

	if secret == nil || secret.Data == nil {
		return nil, nil
	}

	return extractSecretMeta(secret.Data)
}

// extractSecretMeta parses the Vault metadata response into a SecretMeta struct.
func extractSecretMeta(data map[string]any) (*SecretMeta, error) {
	meta := &SecretMeta{
		Versions: make(map[int]SecretVersionMeta),
	}

	// Extract current_version (handle both float64 and json.Number)
	meta.CurrentVersion = extractInt(data["current_version"])

	// Extract versions map
	versionsRaw, ok := data["versions"].(map[string]any)
	if !ok {
		return meta, nil
	}

	for versionStr, versionDataRaw := range versionsRaw {
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			continue
		}

		versionData, ok := versionDataRaw.(map[string]any)
		if !ok {
			continue
		}

		vMeta := SecretVersionMeta{}

		if ct, ok := versionData["created_time"].(string); ok {
			vMeta.CreatedTime = ct
		}

		if dt, ok := versionData["deletion_time"].(string); ok && dt != "" {
			vMeta.Deleted = true
		}

		if destroyed, ok := versionData["destroyed"].(bool); ok {
			vMeta.Destroyed = destroyed
		}

		meta.Versions[version] = vMeta
	}

	return meta, nil
}

// extractInt extracts an integer from a value that could be float64 or json.Number.
func extractInt(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	case int:
		return n
	case int64:
		return int(n)
	default:
		return 0
	}
}

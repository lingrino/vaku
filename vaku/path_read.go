package vaku

import (
	"fmt"
)

// PathRead takes in a PathInput, calls the native vault read on it, extracts the secret,
// and returns it as a map of strings to values. Note that this is the way secrets
// are represented in Vault, as map[string]interface{} (JSON)
func (c *Client) PathRead(i *PathInput) (map[string]interface{}, error) {
	var err error
	var output map[string]interface{}

	// Initialize the input
	i.opType = "read"
	err = c.InitPathInput(i)
	if err != nil {
		return output, fmt.Errorf("failed to init read path %s: %w", i.Path, err)
	}

	// Do the actual read
	secret, err := c.Logical().Read(i.opPath)
	if err != nil {
		return output, fmt.Errorf("failed to read secret at %s: %w", i.opPath, err)
	}
	if secret == nil || secret.Data == nil {
		return output, fmt.Errorf("No value found at %s", i.opPath)
	}

	// V2 Mounts return a nested map[string]interface{} at secret.Data["data"]
	output = secret.Data
	if i.mountVersion == "2" && output != nil {
		metadata := secret.Data["metadata"].(map[string]interface{})
		if metadata["deletion_time"].(string) != "" {
			// Note that path_search and folder_search depend on this VAKU_STATUS
			outputS := map[string]interface{}{
				"VAKU_STATUS": "SECRET_HAS_BEEN_DELETED",
			}
			return outputS, nil
		} else if metadata["destroyed"].(bool) == true {
			outputS := map[string]interface{}{
				"VAKU_STATUS": "SECRET_HAS_BEEN_DESTROYED",
			}
			return outputS, nil
		}


		data := secret.Data["data"]
		if data != nil {
			output = data.(map[string]interface{})
		} else {
			return output, fmt.Errorf("No value found at %s", i.opPath)
		}
	}

	return output, err
}

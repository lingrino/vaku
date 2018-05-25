package vaku

import (
	"fmt"

	"github.com/pkg/errors"
)

// PathRead takes in a path, calls vault read, extracts
// the secret, and returns it as a map of strings to values
func (c *Client) PathRead(i *PathInput) (map[string]interface{}, error) {
	var err error
	var output map[string]interface{}

	// Initialize the input
	i.opType = "read"
	err = c.InitPathInput(i)
	if err != nil {
		return output, errors.Wrapf(err, "Failed to init read path %s", i.Path)
	}

	// Do the actual read
	secret, err := c.Logical().Read(i.opPath)
	if err != nil {
		return output, errors.Wrapf(err, "Failed to read secret at %s", i.opPath)
	}
	if secret == nil || secret.Data == nil {
		return output, fmt.Errorf("No value found at %s", i.opPath)
	}

	// V2 Mounts return a nested map[string]interface{} at secret.Data["data"]
	output = secret.Data
	if i.mountVersion == "2" && output != nil {
		data := secret.Data["data"]
		if data != nil {
			output = data.(map[string]interface{})
		}
	}

	return output, err
}

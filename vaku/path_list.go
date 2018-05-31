package vaku

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

// PathList takes in a PathInput, calls the native vault list on it, extracts
// the secret (list of keys), and returns it. Note that any metadata or other
// information returned by the list is thrown away.
func (c *Client) PathList(i *PathInput) ([]string, error) {
	var err error
	var output []string

	// initialize the input
	i.opType = "list"
	err = c.InitPathInput(i)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize PathInput")
	}

	// do the actual list
	secret, err := c.Logical().List(i.opPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list %s", i.opPath)
	}

	// extract list data from the returned secret
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("Secret at %s was nil", i.opPath)

	}
	data, ok := secret.Data["keys"]
	if !ok || data == nil {
		return nil, fmt.Errorf("No Data[\"keys\"] in secret at %s", i.opPath)
	}
	keys, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Failed to convert keys to interface at %s", i.opPath)
	}

	// Make sure every key is a string and append to output
	for _, k := range keys {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("Failed to assert %s as a string at %s", key, i.opPath)
		}
		output = append(output, key)
	}

	if !i.TrimPathPrefix {
		c.SliceAddKeyPrefix(output, i.Path)
	}

	sort.Strings(output)

	return output, err
}

package vaku

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

// PathList takes in a ListInput, calls vault list, extracts
// the secret, and returns a slice of strings
func (c *Client) PathList(i *ListInput) ([]string, error) {
	var err error
	var output []string

	// initialize the input
	err = c.initListInput(i)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to initialize ListInput")
	}

	// do the actual list
	secret, err := c.client.Logical().List(i.ListPath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list %s", i.ListPath)
	}

	// extract list data from the returned secret
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("Secret at %s was nil", i.ListPath)

	}
	data, ok := secret.Data["keys"]
	if !ok || data == nil {
		return nil, fmt.Errorf("No Data[\"keys\"] in secret at %s", i.ListPath)
	}
	keys, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Failed to convert keys to interface at %s", i.ListPath)
	}

	// Make sure every key is a string and append to output
	for _, k := range keys {
		key, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("Failed to assert %s as a string at %s", key, i.ListPath)
		}
		output = append(output, key)
	}

	if !i.TrimPathPrefix {
		c.SliceAddKeyPrefix(output, i.Path)
	}

	sort.Strings(output)

	return output, err
}

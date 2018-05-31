package vaku

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

// PathSearch takes in a PathInput and a search string, reads the path, and searches
// the read data for a match on the search string. Returns true if the string is found
// in the data. Note that this is a simple search that just checks if the secret
// contains the search string anywhere.
func (c *Client) PathSearch(i *PathInput, s string) (bool, error) {
	var err error

	// Read the data at the path
	read, err := c.PathRead(i)
	if err != nil {
		return false, errors.Wrapf(err, "Failed to read data at path %s", i.Path)
	}

	// We know that read returns map[string]interface{} and that the interface{}
	// is stored by Vault in a json-compatible way. So we can abuse that by marshaling
	// the data into JSON, turning that into a string, and searching the string. This
	// is not the fastest or "right" way to search but it's the least complex and works
	// in this limited space.
	for k, v := range read {
		if strings.Contains(k, s) {
			return true, err
		}
		vjson, err := json.Marshal(v)
		if err != nil {
			return false, errors.Wrapf(err, "failed to marshall value into json for search at path %s", i.Path)
		}
		vstr := string(vjson)
		if strings.Contains(vstr, s) {
			return true, err

		}
	}

	return false, err
}

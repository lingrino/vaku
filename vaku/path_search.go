package vaku

import (
	"encoding/json"
	"fmt"
	"strings"
)

// PathSearch takes in a PathInput and a search string, reads the path, and searches
// the read data for a match on the search string. Returns true if the string is found
// in the data. Note that this is a simple search that just checks if the secret
// contains the search string anywhere. Also note that if this is a KV version 2 mount
// and the input path has been deleted (but not destroyed) this returns false with no error.
func (c *Client) PathSearch(i *PathInput, s string) (bool, error) {
	var err error

	// Read the data at the path
	read, err := c.PathRead(i)
	if err != nil {
		return false, fmt.Errorf("failed to read data at path %s: %w", i.Path, err)
	}

	// We know that read returns map[string]interface{} and that the interface{}
	// is stored by Vault in a json-compatible way. So we can abuse that by marshaling
	// the data into JSON, turning that into a string, and searching the string. This
	// is not the fastest or "right" way to search but it's the least complex and works
	// in this limited space.
	for k, v := range read {
		if strings.Contains(k, "VAKU_STATUS") {
			return false, err
		}
		if strings.Contains(k, s) {
			return true, err
		}
		vjson, err := json.Marshal(v)
		if err != nil {
			return false, fmt.Errorf("failed to marshall value into json for search at path %s: %w", i.Path, err)
		}
		vstr := string(vjson)
		if strings.Contains(vstr, s) {
			return true, err

		}
	}

	return false, err
}

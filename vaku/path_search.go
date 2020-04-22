package vaku

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	// ErrPathSearch when PathSearch fails.
	ErrPathSearch = errors.New("path search")
)

// PathSearch searches for a string at a path. Primitive search that does strings.Contains() on
// string(secret).
func (c *Client) PathSearch(p, s string) (bool, error) {
	read, err := c.PathRead(p)
	if err != nil {
		return false, newWrapErr(p, ErrPathSearch, err)
	}

	match, err := searchSecret(read, s)
	if err != nil {
		return false, newWrapErr(p, ErrPathSearch, err)
	}

	return match, nil
}

// searchSecret does the actual search.
func searchSecret(secret map[string]interface{}, search string) (bool, error) {
	for k, v := range secret {
		if strings.Contains(k, search) {
			return true, nil
		}
		vjson, err := json.Marshal(v)
		if err != nil {
			return false, newWrapErr("", ErrJSONMarshall, nil)
		}
		vstr := string(vjson)
		if strings.Contains(vstr, search) {
			return true, nil
		}
	}

	return false, nil
}

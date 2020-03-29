package vaku

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrPathSearch = errors.New("path search")
)

// PathSearch searches for a string at a path. Primitive search that does strings.Contains() on
// string(secret).
func (c *Client) PathSearch(p, s string) (bool, error) {
	read, err := c.PathRead(p)
	if err != nil {
		return false, newWrapErr(p, ErrPathSearch, err)
	}

	for k, v := range read {
		if strings.Contains(k, s) {
			return true, nil
		}
		vjson, err := json.Marshal(v)
		if err != nil {
			return false, newWrapErr("", ErrPathSearch, ErrJSONMarshall)
		}
		vstr := string(vjson)
		if strings.Contains(vstr, s) {
			return true, nil
		}
	}

	return false, nil
}

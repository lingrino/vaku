package vault

import (
	"fmt"
	"sort"
	"strings"
)

// KVListInput is the required input for KVList
type KVListInput struct {
	Path           string
	Recurse        bool
	TrimPathPrefix bool
}

// NewKVListInput returns a kvListInput struct with default values
func NewKVListInput() *KVListInput {
	return &KVListInput{
		Path:           "",
		Recurse:        false,
		TrimPathPrefix: false,
	}
}

// KVList takes a path and returns a slice of all values at that path
// If Recurse, also list all nested paths/folders
// If TrimPathPrefix, do not prefix keys with leading path
func (c *Client) KVList(i *KVListInput) ([]string, error) {
	var err error
	var result []string

	if i.Path == "" {
		return nil, fmt.Errorf("[FATAL]: KVList: Path is not specified")
	}

	mountPath, version, err := c.PathMountInfo(i.Path)
	if err != nil {
		return nil, fmt.Errorf("[FATAL]: KVList: Failed to describe mount for path %s: %s", i.Path, err)
	}

	if version == "2" {
		i.Path = c.PathJoin(mountPath, "metadata", strings.TrimPrefix(i.Path, mountPath))
	}

	secret, err := c.client.Logical().List(i.Path)
	if err != nil || secret == nil {
		return nil, fmt.Errorf("[FATAL]: KVList: Failed to list path %s: %s", i.Path, err)
	}

	keys, ok := secret.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("[FATAL]: KVList: %s contained no Data['keys']: %s", i.Path, err)
	}

	list, ok := keys.([]interface{})
	if !ok {
		return nil, fmt.Errorf("[FATAL]: KVList: Failed to convert %s keys to interface: %s", i.Path, err)
	}

	for _, v := range list {
		typed, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("[FATAL]: KVList: Failed to assert %s in %s as a string", typed, i.Path)
		}
		if c.PathIsFolder(typed) && i.Recurse {
			subKeys, _ := c.KVList(&KVListInput{
				Path:           c.PathJoin(i.Path, typed),
				Recurse:        i.Recurse,
				TrimPathPrefix: false,
			})
			result = append(result, subKeys...)
		} else if c.PathIsFolder(typed) {
			result = append(result, c.PathJoin(i.Path, typed)+"/")
		} else {
			result = append(result, c.PathJoin(i.Path, typed))
		}
	}

	if i.TrimPathPrefix == true {
		for idx, pth := range result {
			result[idx] = strings.TrimPrefix(strings.TrimPrefix(pth, i.Path), "/")
		}
	}

	sort.Strings(result)

	return result, err
}

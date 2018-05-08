package vault

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// KVListInput is the required input for KVList
type KVListInput struct {
	Path           string
	Recurse        bool
	TrimPathPrefix bool
}

// NewKVListInput takes a path and returns a kvListInput struct with
// default values to produce similar to what is returned by Vault CLI
func NewKVListInput(p string) *KVListInput {
	return &KVListInput{
		Path:           p,
		Recurse:        false,
		TrimPathPrefix: true,
	}
}

// KVList takes a path and returns a slice of all values at that path
// If Recurse, also list all nested paths/folders
// If TrimPathPrefix, do not prefix keys with leading path
func (c *Client) KVList(i *KVListInput) ([]string, error) {
	var err error
	var result []string

	if i.Path == "" {
		return nil, errors.Wrap(err, "Path is not specified")
	}

	mountPath, version, err := c.PathMountInfo(i.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to describe mount for path %s", i.Path)
	}

	if version == "2" {
		i.Path = c.PathJoin(mountPath, "metadata", strings.TrimPrefix(i.Path, mountPath))
	}

	secret, err := c.client.Logical().List(i.Path)
	if err != nil || secret == nil {
		return nil, errors.Wrapf(err, "Failed to list path %s", i.Path)
	}

	keys, ok := secret.Data["keys"]
	if !ok {
		return nil, errors.Wrapf(err, "%s contained no Data['keys']", i.Path)
	}

	list, ok := keys.([]interface{})
	if !ok {
		return nil, errors.Wrapf(err, "Failed to convert %s keys to interface", i.Path)
	}

	for _, v := range list {
		typed, ok := v.(string)
		if !ok {
			return nil, errors.Wrapf(err, "Failed to assert %s in %s as a string", typed, i.Path)
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

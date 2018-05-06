package vault

// import (
// 	"fmt"
// 	"sort"
// )

// // KVListInput is the required input for KVList
// type KVListInput struct {
// 	Path           string
// 	Recurse        bool
// 	TrimPathPrefix bool
// }

// // NewKVListInput returns a kvListInput struct with default values
// func NewKVListInput() *KVListInput {
// 	return &KVListInput{
// 		Path:           "",
// 		Recurse:        false,
// 		TrimPathPrefix: false,
// 	}
// }

// // KVList takes a path and returns a slice of all values at that path
// // If Recurse, also list all nested paths/folders
// // If TrimPathPrefix, do not prefix keys with leading path
// func (c *Client) KVList(i *KVListInput) ([]string, error) {
// 	var err error
// 	var result []string

// 	if i.Path == "" {
// 		return nil, fmt.Errorf("[FATAL]: KVList: Path is not specified")
// 	}

// 	if c.PathIsV2(i.Path) {
// 		i.Path = addPrefixToVKVPath(i.path, mountPath, "metadata")
// 		if err != nil {
// 			c.UI.Error(err.Error())
// 			return nil, 2
// 		}
// 	}

// 	secret, err := c.client.Logical().List(i.Path)
// 	if err != nil || secret == nil {
// 		return nil, fmt.Errorf("[FATAL]: KVList: Failed to list path %s: %s", i.Path, err)
// 	}

// 	keys, ok := secret.Data["keys"]
// 	if !ok {
// 		return nil, fmt.Errorf("[FATAL]: KVList: %s contained no Data['keys']: %s", i.Path, err)
// 	}

// 	list, ok := keys.([]interface{})
// 	if !ok {
// 		return nil, fmt.Errorf("[FATAL]: KVList: Failed to convert %s keys to interface: %s", i.Path, err)
// 	}

// 	for _, v := range list {
// 		typed, ok := v.(string)
// 		if !ok {
// 			return nil, fmt.Errorf("[FATAL]: KVList: Failed to assert %s in %s as a string", typed, i.Path)
// 		}
// 		result = append(result, typed)
// 	}

// 	if i.Recurse {
// 		for _, path := range result {
// 			i.Path = path
// 			subKeys, _ := c.KVList(i)
// 			result = append(result, subKeys...)
// 		}

// 	}

// 	sort.Strings(result)

// 	return result, err
// }

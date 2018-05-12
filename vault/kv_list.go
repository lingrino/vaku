package vault

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// KVListInput is the required input for KVList
type KVListInput struct {
	Path           string
	Recurse        bool
	TrimPathPrefix bool
	MountPath      string
	MountVersion   string
}

// NewKVListInput takes a path and returns a kvListInput struct with
// default values to produce similar to what is returned by Vault CLI
func NewKVListInput(p string) *KVListInput {
	return &KVListInput{
		Path:           p,
		Recurse:        false,
		TrimPathPrefix: true,
		MountPath:      "",
		MountVersion:   "",
	}
}

// listWorkerResult holds the key and any errors from a job
type listWorkerResult struct {
	key string
	err error
}

// listWorker listens on the inputs channel for new paths to list and then does the
// work of listing those paths and returning the result. If recurse is true any listed keys
// that are folders will be added back to the inputs channel for further processing
func (c *Client) listWorker(inputs chan *KVListInput, results chan<- *listWorkerResult, inputsWG *sync.WaitGroup, resultsWG *sync.WaitGroup) {
	// listen on the inputs channel until it is closed
	for {
		i, more := <-inputs
		if more {
			// This lets us only get mount info only once, by passing it in future inputs
			if i.MountPath == "" {
				mountPath, version, err := c.PathMountInfo(i.Path)
				if err != nil {
					fmt.Println("hi")
					resultsWG.Add(1)
					results <- &listWorkerResult{
						key: "",
						err: errors.Wrapf(err, "Failed to describe mount for path %s", i.Path),
					}
					inputsWG.Done()
					continue
				}
				i.MountPath = mountPath
				i.MountVersion = version
			}

			// For v2 mounts lists happen on mount/metadata/path instead of mount/path
			var listPath string
			if i.MountVersion == "2" {
				listPath = c.PathJoin(i.MountPath, "metadata", strings.TrimPrefix(i.Path, i.MountPath))
			} else {
				listPath = i.Path
			}

			// Do the actual list
			secret, err := c.client.Logical().List(listPath)
			if err != nil {
				resultsWG.Add(1)
				results <- &listWorkerResult{
					key: "",
					err: errors.Wrapf(err, "Failed to list path at %s", listPath),
				}
				inputsWG.Done()
				continue
			}

			// extract list data from the returned secret
			if secret == nil || secret.Data == nil {
				resultsWG.Add(1)
				results <- &listWorkerResult{
					key: "",
					err: fmt.Errorf("Secret at %s was nil", listPath),
				}
				inputsWG.Done()
				continue

			}
			keys, ok := secret.Data["keys"]
			if !ok || keys == nil {
				resultsWG.Add(1)
				results <- &listWorkerResult{
					key: "",
					err: fmt.Errorf("No Data[\"keys\"] in secret at %s", listPath),
				}
				inputsWG.Done()
				continue
			}
			list, ok := keys.([]interface{})
			if !ok {
				resultsWG.Add(1)
				results <- &listWorkerResult{
					key: "",
					err: fmt.Errorf("Failed to convert keys to interface at %s", listPath),
				}
				inputsWG.Done()
				continue
			}

			// For each key, either add it to results or add it back to inputs if recurse and folder
			for _, v := range list {
				key, ok := v.(string)
				if !ok {
					resultsWG.Add(1)
					results <- &listWorkerResult{
						key: "",
						err: fmt.Errorf("Failed to assert %s as a string at %s", key, listPath),
					}
					inputsWG.Done()
					continue
				}
				// If we're recursing and the key is a folder, add it back as an input to be listed
				if c.PathIsFolder(key) && i.Recurse {
					inputsWG.Add(1)
					inputs <- &KVListInput{
						Path:           c.PathJoin(i.Path, key),
						Recurse:        i.Recurse,
						TrimPathPrefix: i.TrimPathPrefix,
						MountPath:      i.MountPath,
						MountVersion:   i.MountVersion,
					}
				} else if c.PathIsFolder(key) {
					resultsWG.Add(1)
					results <- &listWorkerResult{
						key: c.PathJoin(i.Path, key) + "/",
						err: nil,
					}

				} else {
					resultsWG.Add(1)
					results <- &listWorkerResult{
						key: c.PathJoin(i.Path, key),
						err: nil,
					}
				}
			}
			inputsWG.Done()
		} else {
			return
		}
	}
}

// KVList takes a path and returns a slice of all values at that path
// If Recurse, also list all nested paths/folders
// If TrimPathPrefix, do not prefix keys with leading path
func (c *Client) KVList(i *KVListInput) ([]string, error) {
	var err error
	var output []string

	if i.Path == "" {
		return nil, errors.Wrap(err, "Path is not specified")
	}

	inputs := make(chan *KVListInput, 5)
	results := make(chan *listWorkerResult, 5)

	var inputsWG sync.WaitGroup
	var resultsWG sync.WaitGroup

	// Add our first input
	inputsWG.Add(1)
	inputs <- i

	// Listen on results channel and add keys to an output list
	go func() {
		for {
			o, more := <-results
			if more {
				if o.err != nil {
					err = errors.Wrapf(err, "Failed to list path %s", i.Path)
				} else {
					output = append(output, o.key)
				}

				resultsWG.Done()
			} else {
				return
			}
		}
	}()

	// Spawn 5 workers if recursing, otherise just one
	// TODO - read worker count from configuration
	if i.Recurse {
		for w := 1; w <= 5; w++ {
			go c.listWorker(inputs, results, &inputsWG, &resultsWG)
		}
	} else {
		go c.listWorker(inputs, results, &inputsWG, &resultsWG)
	}

	// Wait until all lists are complete
	inputsWG.Wait()
	resultsWG.Wait()
	close(inputs)
	close(results)

	// Remove the prefix if it is not wanted
	if i.TrimPathPrefix == true {
		for idx, pth := range output {
			output[idx] = strings.TrimPrefix(strings.TrimPrefix(pth, i.Path), "/")
		}
	}

	sort.Strings(output)

	return output, err
}

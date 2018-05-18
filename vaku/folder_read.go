package vaku

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// folderReadWorkerOutput holds the key and any errors from a job
type folderReadWorkerOutput struct {
	readPath string
	data     map[string]interface{}
	err      error
}

// folderReadWorkerInput takes input/output channels and
// waitgroups to update as new work is discovered
type folderReadWorkerInput struct {
	inputsC  chan *PathInput
	resultsC chan<- *folderReadWorkerOutput
}

func (c *Client) folderReadWorker(i *folderReadWorkerInput) {
	for {
		pi, more := <-i.inputsC
		if more {
			read, err := c.PathRead(pi)
			if err != nil {
				i.resultsC <- &folderReadWorkerOutput{
					readPath: "",
					data:     nil,
					err:      errors.Wrapf(err, "Failed to read path %s", pi.Path),
				}
				fmt.Println(err)
				continue
			}
			i.resultsC <- &folderReadWorkerOutput{
				readPath: pi.Path,
				data:     read,
				err:      nil,
			}
		} else {
			return
		}
	}
}

// folderReadCaller does the actual work of scheduling the reads and collecting the
// results, since that work is shared for FolderRead and FolderReadAll
func (c *Client) folderReadCaller(i *PathInput, keys []string) (map[string]map[string]interface{}, error) {
	var err error
	var output map[string]map[string]interface{}

	// Initialize the input
	i.opType = "read"
	c.InitPathInput(i)

	// Remove folders (can't be read) from the list
	keys = c.SliceRemoveFolders(keys)

	// Concurrency channels for workers
	// Create output equal to length of keys
	inputsC := make(chan *PathInput, len(keys))
	resultsC := make(chan *folderReadWorkerOutput, len(keys))
	output = make(map[string]map[string]interface{}, len(keys))

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderReadWorker(&folderReadWorkerInput{
			inputsC:  inputsC,
			resultsC: resultsC,
		})
	}

	// Add all paths to read to the inputs channel
	for _, p := range keys {
		inputsC <- &PathInput{
			Path:           p,
			opType:         i.opType,
			mountPath:      i.mountPath,
			mountlessPath:  i.mountlessPath,
			mountVersion:   i.mountVersion,
			TrimPathPrefix: false,
		}
	}
	close(inputsC)

	// Empty the results channel into output
	for j := 0; j < len(keys); j++ {
		o := <-resultsC
		if o.err != nil {
			err = errors.Wrapf(o.err, "Failed to read path %s", o.readPath)
		} else {
			if i.TrimPathPrefix {
				output[c.KeyJoin(strings.TrimPrefix(o.readPath, i.Path))] = o.data
			} else {
				output[o.readPath] = o.data
			}
		}
	}

	return output, err
}

// FolderRead takes in a PathInput, reads all non-folders in that path
// and outputs a map of paths to values at that path
func (c *Client) FolderRead(i *PathInput) (map[string]map[string]interface{}, error) {
	var err error
	var output map[string]map[string]interface{}

	// Get the keys to read
	list, err := c.PathList(&PathInput{
		Path:           i.Path,
		TrimPathPrefix: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list %s", i.Path)
	}

	// Hand over to folderReadCaller
	output, err = c.folderReadCaller(i, list)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read folder at %s", i.Path)
	}

	return output, err
}

// FolderReadAll takes in a PathInput, reads all keys in that path
// and all nested paths and outputs a map of paths to values at that path
func (c *Client) FolderReadAll(i *PathInput) (map[string]map[string]interface{}, error) {
	var err error
	var output map[string]map[string]interface{}

	// Get the keys to read
	list, err := c.FolderList(&PathInput{
		Path:           i.Path,
		TrimPathPrefix: false,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list %s", i.Path)
	}

	// Hand over to folderReadCaller
	output, err = c.folderReadCaller(i, list)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read folder at %s", i.Path)
	}

	return output, err
}

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

// FolderRead takes in a PathInput, reads all non-folders in that path
// and outputs a map of paths to values at that path
// TODO - add ability to recurse before read
func (c *Client) FolderRead(i *PathInput) (map[string]map[string]interface{}, error) {
	var err error
	var output map[string]map[string]interface{}

	// Don't trim prefix during indivudal reads, only at end
	trimPrefix := i.TrimPathPrefix
	i.TrimPathPrefix = false

	// Get the keys to read
	list, err := c.PathList(i)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to list %s", i.Path)
	}

	// Remove folders (can't be read) from the list
	keys := c.SliceRemoveFolders(list)

	// Concurrency channels for workers
	// Create output equal to length of keys
	inputsC := make(chan *PathInput, len(keys))
	resultsC := make(chan *folderReadWorkerOutput, len(keys))
	output = make(map[string]map[string]interface{}, len(keys))

	// Spawn 5 workers
	// TODO - read worker/concurrency count from configuration
	for w := 1; w <= 5; w++ {
		go c.folderReadWorker(&folderReadWorkerInput{
			inputsC:  inputsC,
			resultsC: resultsC,
		})
	}

	// Add all paths to read to the inputs channel
	for _, p := range keys {
		inputsC <- NewPathInput(p)
	}
	close(inputsC)

	// Empty the results channel into output
	for j := 0; j < len(keys); j++ {
		o := <-resultsC
		if o.err != nil {
			err = errors.Wrapf(o.err, "Failed to read path %s", o.readPath)
		} else {
			if trimPrefix {
				output[c.KeyJoin(strings.TrimPrefix(o.readPath, i.Path))] = o.data
			} else {
				output[o.readPath] = o.data
			}
		}

	}

	return output, err
}

package vaku

import (
	"github.com/pkg/errors"
)

// folderWriteWorkerInput takes input/output channels for input to the job
type folderWriteWorkerInput struct {
	inputsC  <-chan *writeInput
	resultsC chan<- error
}

type writeInput struct {
	path *PathInput
	data map[string]interface{}
}

// FolderWrite takes in a map of paths to data that should be written to that path. Users may
// find this function difficult to call on its own because the input can be large and specific,
// however the output of FolderRead matches the input of FolderWrite(), making them easy to use
// together. Note that mount/version information is determined only once using a random path in
// the map and cached for all future writes. Therefore this function cannot write to two mounts of
// different versions in the same call.
func (c *Client) FolderWrite(d map[string]map[string]interface{}) error {
	var err error
	var basePathInfo *PathInput

	// Get mount data based on a random path in the map
	for p := range d {
		basePathInfo = NewPathInput(p)
		break
	}
	basePathInfo.opType = "write"
	c.InitPathInput(basePathInfo)

	// Concurrency channels for workers
	inputsC := make(chan *writeInput, len(d))
	resultsC := make(chan error, len(d))

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderWriteWorker(&folderWriteWorkerInput{
			inputsC:  inputsC,
			resultsC: resultsC,
		})
	}

	// Add all path/data parts to write to the inputs channel
	for k, v := range d {
		inputsC <- &writeInput{
			path: &PathInput{
				Path:          k,
				mountPath:     basePathInfo.mountPath,
				mountlessPath: basePathInfo.mountlessPath,
				mountVersion:  basePathInfo.mountVersion,
			},
			data: v,
		}
	}
	close(inputsC)

	// Empty the results channel into output
	for j := 0; j < len(d); j++ {
		o := <-resultsC
		if o != nil {
			err = errors.Wrap(o, "Failed to write path")
		}
	}

	return err
}

func (c *Client) folderWriteWorker(i *folderWriteWorkerInput) {
	var err error
	for {
		id, more := <-i.inputsC
		if more {
			err = c.PathWrite(id.path, id.data)
			if err != nil {
				i.resultsC <- errors.Wrapf(err, "Failed to write path %s", id.path)
				continue
			}
			i.resultsC <- nil
		} else {
			return
		}
	}
}

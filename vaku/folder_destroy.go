package vaku

import (
	"fmt"
)

// folderDestroyWorkerInput takes input/output channels for input to the job
type folderDestroyWorkerInput struct {
	inputsC  <-chan *PathInput
	resultsC chan<- error
}

// FolderDestroy takes in a path and destroys every key in that folder and all sub-folders.
// Note that this function only works on V2 mounts and that it destroys ALL versions of ALL keys
func (c *Client) FolderDestroy(i *PathInput) error {
	var err error

	// Get the keys to destroy
	list, err := c.FolderList(&PathInput{
		Path:           i.Path,
		TrimPathPrefix: false,
	})
	if err != nil {
		return fmt.Errorf("failed to list %s: %w", i.Path, err)
	}

	// Init the path
	i.opType = "destroy"
	err = c.InitPathInput(i)
	if err != nil {
		return fmt.Errorf("failed to init path %s: %w", i.Path, err)
	}

	// Concurrency channels for workers
	inputsC := make(chan *PathInput, len(list))
	resultsC := make(chan error, len(list))

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderDestroyWorker(&folderDestroyWorkerInput{
			inputsC:  inputsC,
			resultsC: resultsC,
		})
	}

	// Add all paths to destroy to the inputs channel
	for _, p := range list {
		inputsC <- &PathInput{
			Path:          p,
			mountPath:     i.mountPath,
			mountlessPath: i.mountlessPath,
			mountVersion:  i.mountVersion,
		}
	}
	close(inputsC)

	// Empty the results channel into output
	for j := 0; j < len(list); j++ {
		o := <-resultsC
		if o != nil {
			err = fmt.Errorf("failed to destroy path: %w", o)
		}
	}

	return err
}

// folderDestroyWorker does the work of reading a path from a channel and destroying it
func (c *Client) folderDestroyWorker(i *folderDestroyWorkerInput) {
	var err error
	for {
		path, more := <-i.inputsC
		if more {
			err = c.PathDestroy(path)
			if err != nil {
				i.resultsC <- fmt.Errorf("failed to destroy path %s: %w", path.Path, err)
				continue
			}
			i.resultsC <- nil
		} else {
			return
		}
	}
}

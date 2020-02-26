package vaku

import (
	"fmt"
)

// folderDeleteWorkerInput takes input/output channels for input to the job
type folderDeleteWorkerInput struct {
	inputsC  <-chan *PathInput
	resultsC chan<- error
}

// FolderDelete takes in a path and deletes every key in that folder and all sub-folders.
// Note that this calls PathDelete() on every path found in the folder, and for v2 secret
// mounts that means deleting the active version, but not all versions.
func (c *Client) FolderDelete(i *PathInput) error {
	var err error

	// Get the keys to delete
	list, err := c.FolderList(&PathInput{
		Path:           i.Path,
		TrimPathPrefix: false,
	})
	if err != nil {
		return fmt.Errorf("failed to list %s: %w", i.Path, err)
	}

	// Init the path
	i.opType = "delete"
	err = c.InitPathInput(i)
	if err != nil {
		return fmt.Errorf("failed to init path %s: %w", i.Path, err)
	}

	// Concurrency channels for workers
	inputsC := make(chan *PathInput, len(list))
	resultsC := make(chan error, len(list))

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderDeleteWorker(&folderDeleteWorkerInput{
			inputsC:  inputsC,
			resultsC: resultsC,
		})
	}

	// Add all paths to delete to the inputs channel
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
			err = fmt.Errorf("failed to delete path: %w", o)
		}
	}

	return err
}

// folderDeleteWorker does the work of reading a path from a channel and deleting it
func (c *Client) folderDeleteWorker(i *folderDeleteWorkerInput) {
	var err error

	for {
		path, more := <-i.inputsC
		if more {
			err = c.PathDelete(path)
			if err != nil {
				i.resultsC <- fmt.Errorf("failed to delete path %s: %w", path.Path, err)
				continue
			}
			i.resultsC <- nil
		} else {
			return
		}
	}
}

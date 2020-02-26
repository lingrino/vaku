package vaku

import (
	"fmt"
	"sort"
	"strings"
)

// folderSearchWorkerInput takes input/output channels for input to the job
type folderSearchWorkerInput struct {
	inputsC  <-chan *PathInput
	searchS  string
	resultsC chan<- *folderSearchWorkerOutput
}

// folderSearchWorkerOutput is the data returned by a FolderSearch worker
type folderSearchWorkerOutput struct {
	path  string
	match bool
	err   error
}

// FolderSearch takes in a PathInput and a search string. It then calls FolderList() on that path
// and concurrently runs PathSearch() on every returned path in the list. It returns a list of paths
// at which the search string was found. Note that running this function against a path that is not
// a folder will produce an error, to search a single path you should use PathSearch() instead.
func (c *Client) FolderSearch(i *PathInput, s string) ([]string, error) {
	var err error
	var output []string

	// Init the path to get mount info
	i.opType = "read"
	err = c.InitPathInput(i)
	if err != nil {
		return output, fmt.Errorf("failed to init path %s: %w", i.Path, err)
	}

	// Get all of the paths to search
	list, err := c.FolderList(&PathInput{
		Path:           i.Path,
		TrimPathPrefix: true,
	})
	if err != nil {
		return output, fmt.Errorf("failed to list folder at %s: %w", i.Path, err)
	}

	// Concurrency channels for workers
	inputsC := make(chan *PathInput, len(list))
	resultsC := make(chan *folderSearchWorkerOutput, len(list))

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderSearchWorker(&folderSearchWorkerInput{
			inputsC:  inputsC,
			searchS:  s,
			resultsC: resultsC,
		})
	}

	// Add all paths to search to the path channel
	for _, p := range list {
		inputsC <- &PathInput{
			Path:          c.PathJoin(i.Path, p),
			mountPath:     i.mountPath,
			mountlessPath: i.mountlessPath,
			mountVersion:  i.mountVersion,
		}
	}
	close(inputsC)

	// Empty the results channel into output
	for j := 0; j < len(list); j++ {
		o := <-resultsC
		if o.err != nil {
			err = fmt.Errorf("failed to search path %s: %w", o.path, err)
		}
		if o.match {
			if i.TrimPathPrefix {
				outputPath := c.PathJoin(strings.TrimPrefix(o.path, i.Path))
				output = append(output, outputPath)
			} else {
				output = append(output, o.path)
			}
		}
	}

	sort.Strings(output)

	return output, err
}

// folderSearchWorker does the work of searching a single path
func (c *Client) folderSearchWorker(i *folderSearchWorkerInput) {
	for {
		input, more := <-i.inputsC
		if more {
			match, err := c.PathSearch(input, i.searchS)
			output := &folderSearchWorkerOutput{
				err:   err,
				path:  input.Path,
				match: false,
			}
			if err != nil {
				output.err = fmt.Errorf("failed to search path %s: %w", input.Path, err)
				i.resultsC <- output
				continue
			}
			if match {
				output.match = true
			}
			i.resultsC <- output
		} else {
			return
		}
	}
}

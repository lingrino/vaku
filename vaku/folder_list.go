package vaku

import (
	"sort"
	"sync"

	"github.com/pkg/errors"
)

// folderListWorkerOutput holds the key and any errors from a job
type folderListWorkerOutput struct {
	key string
	err error
}

// folderListWorkerInput takes input/output channels and
// waitgroups to update as new work is discovered
type folderListWorkerInput struct {
	inputsC   chan *PathInput
	resultsC  chan<- *folderListWorkerOutput
	inputsWG  *sync.WaitGroup
	resultsWG *sync.WaitGroup
}

// FolderList takes in a PathInput and walks the path by calling PathList
// on the input path and all folders within that path as well. Returns the
// results as a sorted slice of paths.
func (c *Client) FolderList(i *PathInput) ([]string, error) {
	var err error
	var output []string

	// Don't trim prefix during indivudal lists, only at end
	trimPrefix := i.TrimPathPrefix
	i.TrimPathPrefix = false

	// Concurrency tools for waiting on workers
	inputsC := make(chan *PathInput, 5)
	resultsC := make(chan *folderListWorkerOutput, 5)
	var inputsWG sync.WaitGroup
	var resultsWG sync.WaitGroup

	// Add our first input
	inputsWG.Add(1)
	inputsC <- i

	// Listen on results channel and add keys to an output list
	go func() {
		for {
			o, more := <-resultsC
			if more {
				if o.err != nil {
					err = errors.Wrapf(o.err, "Failed to list path %s", i.Path)
				} else {
					output = append(output, o.key)
				}
				resultsWG.Done()
			} else {
				return
			}
		}
	}()

	// Spawn workers equal to MaxConcurrency
	for w := 1; w <= MaxConcurrency; w++ {
		go c.folderListWorker(&folderListWorkerInput{
			inputsC:   inputsC,
			resultsC:  resultsC,
			inputsWG:  &inputsWG,
			resultsWG: &resultsWG,
		})
	}

	// Wait until everything is listed and tear down
	inputsWG.Wait()
	resultsWG.Wait()
	close(inputsC)
	close(resultsC)

	if trimPrefix {
		c.SliceTrimKeyPrefix(output, i.Path)
	}

	sort.Strings(output)

	return output, err
}

func (c *Client) folderListWorker(i *folderListWorkerInput) {
	for {
		l, more := <-i.inputsC
		if more {
			list, err := c.PathList(l)
			if err != nil {
				i.resultsWG.Add(1)
				i.resultsC <- &folderListWorkerOutput{
					key: "",
					err: errors.Wrapf(err, "Failed to list path %s", l.Path),
				}
				i.inputsWG.Done()
				continue
			}
			for _, key := range list {
				if c.KeyIsFolder(key) {
					i.inputsWG.Add(1)
					i.inputsC <- &PathInput{
						Path:           c.PathJoin(l.Path, c.KeyBase(key)),
						opPath:         c.PathJoin(l.opPath, c.KeyBase(key)),
						mountPath:      l.mountPath,
						mountVersion:   l.mountVersion,
						TrimPathPrefix: l.TrimPathPrefix,
					}
				} else {
					i.resultsWG.Add(1)
					i.resultsC <- &folderListWorkerOutput{
						key: key,
						err: nil,
					}
				}
			}
			i.inputsWG.Done()
		} else {
			return
		}
	}
}

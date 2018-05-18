package vaku

import (
	"fmt"

	"github.com/pkg/errors"
)

// PathInput is the input for List
type PathInput struct {
	Path           string
	TrimPathPrefix bool
	opPath         string
	opType         string
	mountPath      string
	mountVersion   string
}

// NewPathInput takes in a Path and returns
// the default PathInput. Only Path is required
func NewPathInput(p string) *PathInput {
	return &PathInput{
		Path:           p,
		TrimPathPrefix: true,
		opPath:         "",
		opType:         "",
		mountPath:      "",
		mountVersion:   "",
	}
}

// InitPathInput fills in missing values from PathInput
// with defaults and mount information
func (c *Client) InitPathInput(i *PathInput) error {
	var err error

	if i.Path == "" {
		return fmt.Errorf("Path is required and not specified")
	}

	if i.opPath == "" || i.mountPath == "" || i.mountVersion == "" {
		m, err := c.MountInfo(i.Path)
		if err != nil {
			return errors.Wrapf(err, "Failed to describe mount for path %s", i.Path)
		}

		if m.mountVersion == "2" {
			if i.opType == "list" {
				i.opPath = c.PathJoin(m.mountPath, "metadata", m.MountlessPath)
			} else if i.opType == "read" {
				i.opPath = c.PathJoin(m.mountPath, "data", m.MountlessPath)
			}
		} else {
			i.opPath = c.PathJoin(i.Path)
		}

		i.mountPath = m.mountPath
		i.mountVersion = m.mountVersion
	}

	return err
}

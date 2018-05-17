package vaku

import (
	"fmt"

	"github.com/pkg/errors"
)

// PathInput is the input for List
type PathInput struct {
	Path           string
	OpPath         string
	OpType         string
	MountPath      string
	MountVersion   string
	TrimPathPrefix bool
}

// NewPathInput takes in a Path and returns
// the default PathInput. Only Path is required
func NewPathInput(p string) *PathInput {
	return &PathInput{
		Path:           p,
		OpPath:         "",
		OpType:         "",
		MountPath:      "",
		MountVersion:   "",
		TrimPathPrefix: true,
	}
}

// InitPathInput fills in missing values from PathInput
// with defaults and mount information
func (c *Client) InitPathInput(i *PathInput) error {
	var err error

	if i.Path == "" {
		return fmt.Errorf("Path is required and not specified")
	}

	if i.OpPath == "" || i.MountPath == "" || i.MountVersion == "" {
		m, err := c.MountInfo(i.Path)
		if err != nil {
			return errors.Wrapf(err, "Failed to describe mount for path %s", i.Path)
		}

		if m.MountVersion == "2" {
			if i.OpType == "list" {
				i.OpPath = c.PathJoin(m.MountPath, "metadata", m.MountlessPath)
			} else if i.OpType == "read" {
				i.OpPath = c.PathJoin(m.MountPath, "data", m.MountlessPath)
			}
		} else {
			i.OpPath = c.PathJoin(i.Path)
		}

		i.MountPath = m.MountPath
		i.MountVersion = m.MountVersion
	}

	return err
}

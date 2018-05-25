package vaku

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// PathInput is the input for List
type PathInput struct {
	Path           string
	TrimPathPrefix bool
	opPath         string
	opType         string
	mountPath      string
	mountlessPath  string
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
		mountlessPath:  "",
		mountVersion:   "",
	}
}

// InitPathInput fills in missing values from PathInput
// with defaults and mount information
func (c *Client) InitPathInput(i *PathInput) error {
	var err error

	// Required values
	if i.Path == "" {
		return fmt.Errorf("Path is required and not specified")
	}
	if i.opType == "" {
		return fmt.Errorf("opType is required and not specified")
	}

	// If mount info is already set don't get again. Only ensure that i.opPath is correct
	// Otherwise populate based on MountInfo for the path
	if i.mountPath != "" && i.mountVersion != "" {
		i.mountlessPath = strings.TrimPrefix(i.Path, i.mountPath)
		if i.mountVersion == "2" {
			if i.opType == "list" {
				i.opPath = c.PathJoin(i.mountPath, "metadata", i.mountlessPath)
			} else if i.opType == "read" || i.opType == "write" {
				i.opPath = c.PathJoin(i.mountPath, "data", i.mountlessPath)
			}
		} else {
			i.opPath = c.PathJoin(i.Path)
		}
	} else if i.opPath == "" || i.mountPath == "" || i.mountVersion == "" || i.mountlessPath == "" {
		m, err := c.MountInfo(i.Path)
		if err != nil {
			return errors.Wrapf(err, "Failed to describe mount for path %s", i.Path)
		}
		if m.MountVersion == "2" {
			if i.opType == "list" {
				i.opPath = c.PathJoin(m.MountPath, "metadata", m.MountlessPath)
			} else if i.opType == "read" || i.opType == "write" {
				i.opPath = c.PathJoin(m.MountPath, "data", m.MountlessPath)
			}
		} else {
			i.opPath = c.PathJoin(i.Path)
		}

		i.mountPath = m.MountPath
		i.mountVersion = m.MountVersion
		i.mountlessPath = m.MountlessPath
	}

	return err
}

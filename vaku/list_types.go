package vaku

import (
	"fmt"

	"github.com/pkg/errors"
)

// ListInput is the input for List
type ListInput struct {
	Path           string
	ListPath       string
	MountPath      string
	MountVersion   string
	TrimPathPrefix bool
}

// NewListInput takes in a Path and returns
// the default ListInput. Only Path is required
func NewListInput(p string) *ListInput {
	return &ListInput{
		Path:           p,
		ListPath:       "",
		MountPath:      "",
		MountVersion:   "",
		TrimPathPrefix: true,
	}
}

// initListInput fills in missing values from ListInput
// with defaults and mount information
func (c *Client) initListInput(i *ListInput) error {
	var err error

	if i.Path == "" {
		return fmt.Errorf("Path is required and not specified")
	}

	if i.ListPath == "" || i.MountPath == "" || i.MountVersion == "" {
		m, err := c.MountInfo(i.Path)
		if err != nil {
			return errors.Wrapf(err, "Failed to describe mount for path %s", i.Path)
		}

		if m.MountVersion == "2" {
			i.ListPath = c.PathJoin(m.MountPath, "metadata", m.MountlessPath)
		} else {
			i.ListPath = c.PathJoin(i.Path)
		}

		i.MountPath = m.MountPath
		i.MountVersion = m.MountVersion
	}

	return err
}

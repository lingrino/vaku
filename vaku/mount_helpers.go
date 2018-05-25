package vaku

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// MountInfoOutput holds output for MountInfo
// FullPath is the original input path
// mountPath is the path of the mount
// MountlessPath is the FullPath with the mountPath removed
// mountVersion is the version of the mount, or "unknown"
type MountInfoOutput struct {
	FullPath      string
	MountPath     string
	MountlessPath string
	MountVersion  string
}

// MountInfo gets information for the mount of the specified path
func (c *Client) MountInfo(p string) (*MountInfoOutput, error) {
	var err error
	var output MountInfoOutput

	mounts, err := c.Sys().ListMounts()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list mounts")
	}

	// Assumes that all mounts have unique prefixes
	// i.e that there cannot a sec/ and secret/
	for mount, data := range mounts {
		if strings.HasPrefix(p, mount) {
			version, ok := data.Options["version"]
			if !ok {
				version = "unknown"
			}

			output.FullPath = c.PathJoin(p)
			output.MountPath = c.PathJoin(mount)
			output.MountlessPath = c.PathJoin(strings.TrimPrefix(p, mount))
			output.MountVersion = version
			continue
		}
	}
	// Use this as a check to see if a mount was found
	if output.FullPath == "" {
		return nil, fmt.Errorf("MountInfo: Mount not found")
	}

	return &output, err
}

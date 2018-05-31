package vaku

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// MountInfoOutput holds output for MountInfo. This data can be useful for
// determining which key/value engine version a path is mounted on and acting
// accordingly.
type MountInfoOutput struct {
	FullPath      string // The original input path
	MountPath     string // The path of the mount
	MountlessPath string // The FullPath with the MountPath removed
	MountVersion  string // The version of the mount, default "unknown"
}

// MountInfo gets information about the mount at the specified path that can be
// used to determine what actions to take on that path.
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

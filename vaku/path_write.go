package vaku

import (
	"github.com/pkg/errors"
)

// PathWrite takes in a PathInput and data to written to that path. It then
// calls the native vault write with that data at the specified path.
func (c *Client) PathWrite(i *PathInput, d map[string]interface{}) error {
	var err error

	// Initialize the input
	i.opType = "write"
	err = c.InitPathInput(i)
	if err != nil {
		return errors.Wrapf(err, "Failed to init write path %s", i.Path)
	}

	// V2 mounts nest the actual data in another map[string]interface{}
	// https://github.com/hashicorp/vault/blob/69b1cae9e252e9f2f8394675f8df5cd9dca8f5de/command/kv_put.go#L130-L142
	if i.mountVersion == "2" {
		d = map[string]interface{}{
			"data": d,
		}
	}

	// Do the actual write
	_, err = c.Logical().Write(i.opPath, d)
	if err != nil {
		return errors.Wrapf(err, "Failed to write secret to %s", i.opPath)
	}

	return err
}

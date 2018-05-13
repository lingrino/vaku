package vault

import (
	vapi "github.com/hashicorp/vault/api"
	"github.com/pkg/errors"
)

// Client is a wrapper around a real Vault API client.
type Client struct {
	client *vapi.Client
}

// NewClient Returns a new empty Client type
func NewClient() *Client {
	return &Client{}
}

// simpleInit initializes a new default Vault Client
// it should only be used internally for testing
func (c *Client) simpleInit() error {
	var err error

	client, err := vapi.NewClient(vapi.DefaultConfig())
	if err != nil {
		return errors.Wrap(err, "Failed to init the vault client")
	}
	c.client = client

	return err
}

// seed uses a client to write dummy data used for testing to vault
// strings generated here: https://www.random.org/strings
func (c *Client) seed() error {
	var err error
	var mountPath string

	c.client.Sys().EnableAuditWithOptions("audit_stdout", &vapi.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": "stdout",
			"log_raw":   "true",
		},
	})

	seeds := map[string]map[string]string{
		"test/foo": {
			"value": "bar",
		},
		"test/value": {
			"fizz": "buzz",
			"foo":  "bar",
		},
		"test/fizz": {
			"fizz": "buzz",
			"foo":  "bar",
		},
		"test/HToOeKKD": {
			"3zqxVbJY": "TvOjGxvC",
		},
		"test/inner/WKNC3muM": {
			"IY1C148K": "JxBfEt91",
			"iwVzPqbY": "0NH9GlR1",
		},
		"test/inner/A2xlzTfE": {
			"Eg5ljS7t": "BHRMKjj1",
			"quqr32S5": "pcidzSMW",
		},
		"test/inner/again/inner/UCrt6sZT": {
			"Eg5ljS7t": "6F1B5nBg",
			"quqr32S5": "81iY4HAN",
			"r6R0JUzX": "rs1mCRB5",
		},
	}

	// Seed v1 mount
	mountPath = "secretv1/"
	err = c.client.Sys().Mount(mountPath, &vapi.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "1",
		},
	})
	for path, secret := range seeds {
		writePath := mountPath + path
		data := make(map[string]interface{})

		for k, v := range secret {
			data[k] = v
		}

		_, err = c.client.Logical().Write(writePath, data)
		if err != nil {
			return errors.Wrapf(err, "Failed to seed vault at path %s", writePath)
		}
	}

	// Seed v2 mount
	mountPath = "secretv2/"
	c.client.Sys().Mount(mountPath, &vapi.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	})
	for path, secret := range seeds {
		writePath := mountPath + "data/" + path
		data := make(map[string]interface{})

		for k, v := range secret {
			data[k] = v
		}

		// For v2 API
		// https://github.com/hashicorp/vault/blob/master/command/kv_put.go#L130-L142
		data = map[string]interface{}{
			"data": data,
		}

		_, err = c.client.Logical().Write(writePath, data)
		if err != nil {
			return errors.Wrapf(err, "Failed to seed vault at path %s", writePath)
		}
	}
	return err
}

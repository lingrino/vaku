package vault

import (
	"fmt"

	vapi "github.com/hashicorp/vault/api"
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
		return fmt.Errorf("[FATAL]: simpleInit: Failed to init the vault client: %s", err)
	}
	c.client = client

	return err
}

// seed uses a client to write dummy data used for teting to vault
// strings generated here: https://www.random.org/strings
func (c *Client) seed() error {
	var err error

	seeds := map[string]map[string]string{
		"secret/data/test/foo": {
			"value": "bar",
		},
		"secret/data/test/value": {
			"fizz": "buzz",
			"foo":  "bar",
		},
		"secret/data/test/fizz": {
			"fizz": "buzz",
			"foo":  "bar",
		},
		"secret/data/test/HToOeKKD": {
			"3zqxVbJY": "TvOjGxvC",
		},
		"secret/data/test/inner/WKNC3muM": {
			"IY1C148K": "JxBfEt91",
			"iwVzPqbY": "0NH9GlR1",
		},
		"secret/data/test/inner/A2xlzTfE": {
			"Eg5ljS7t": "BHRMKjj1",
			"quqr32S5": "pcidzSMW",
		},
		"secret/data/test/inner/again/inner/UCrt6sZT": {
			"Eg5ljS7t": "6F1B5nBg",
			"quqr32S5": "81iY4HAN",
			"r6R0JUzX": "rs1mCRB5",
		},
	}

	for path, secret := range seeds {
		data := make(map[string]interface{})

		for k, v := range secret {
			data[k] = v
		}

		// For v2 API
		// https://github.com/hashicorp/vault/blob/master/command/kv_put.go#L130-L142
		data = map[string]interface{}{
			"data": data,
		}

		_, err = c.client.Logical().Write(path, data)
		if err != nil {
			return fmt.Errorf("[FATAL]: seed: Failed to seed vault at path %s: %s", path, err)
		}
	}
	return err
}

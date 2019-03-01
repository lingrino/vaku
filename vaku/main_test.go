package vaku_test

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/lingrino/vaku/vaku"
	"github.com/pkg/errors"
)

var seededOnce = false

// Initialize a new simple vault client to be used for tets
func clientInitForTests(t *testing.T) *vaku.Client {
	// Initialize a new vault client
	vclient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal(errors.Wrapf(err, "Failed to create a vault client for testing"))
	}

	// Initialize a new vaku client and attach the vault client
	client := vaku.NewClient()
	client.Client = vclient

	// Set the address and token to the test values
	client.SetToken(vaultToken)
	client.SetAddress(vaultAddr)

	// Seed the client if it has never been seeded
	if !seededOnce {
		err = seed(t, client)
		if err != nil {
			t.Fatal(errors.Wrapf(err, "Failed to seed the vault client"))
		}
		seededOnce = true
	}
	return client
}

// seed uses a client to write dummy data used for testing to vault.
// Strings generated here: https://www.random.org/strings
func seed(t *testing.T, c *vaku.Client) error {
	t.Helper()
	var err error

	// Turn on logging to stdout
	c.Sys().EnableAuditWithOptions("audit_stdout", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": "stdout",
			"log_raw":   "true",
		},
	})

	// Mount the two secret backends
	c.Sys().Mount("secretv1/", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "1"},
	})
	c.Sys().Mount("secretv2/", &api.MountInput{
		Type: "kv",
		Options: map[string]string{
			"version": "2",
		},
	})

	seeds := map[string]map[string]interface{}{
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

	v1Seeds := make(map[string]map[string]interface{})
	v2Seeds := make(map[string]map[string]interface{})
	for k, v := range seeds {
		v1Seeds[c.PathJoin("secretv1", k)] = v
		v2Seeds[c.PathJoin("secretv2", k)] = v
	}

	err = c.FolderWrite(v1Seeds)
	if err != nil {
		return errors.Wrap(err, "Failed to seed secretv1 path")
	}
	err = c.FolderWrite(v2Seeds)
	if err != nil {
		return errors.Wrap(err, "Failed to seed secretv2 path")
	}

	return err
}

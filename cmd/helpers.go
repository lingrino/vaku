package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Lingrino/vaku/vaku"
	"github.com/pkg/errors"

	vapi "github.com/hashicorp/vault/api"
)

// vgc is the vaku client used by CLI commands
var vgc *vaku.Client

// authVGC initializes the vgc, vakuGlobalClient, to be used by all CLI commands
// TODO - try read from ~/.vault-token
func authVGC() {
	// Initialize a new vault client
	vclient, err := vapi.NewClient(vapi.DefaultConfig())
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to create vault client"))
		os.Exit(1)
	} else if vclient.Token() == "" {
		fmt.Println(errors.New("VAULT_TOKEN not set"))
		os.Exit(1)
	}

	// Add the Vault client to the Vaku client
	vgc = vaku.NewClient()
	vgc.Client = vclient
}

func print(i map[string]interface{}) {
	json, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}

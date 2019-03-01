package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/lingrino/vaku/vaku"
	vapi "github.com/hashicorp/vault/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// vgc is the vaku client used by CLI commands
var vgc *vaku.Client

// authVGC initializes the vgc (vaku global client) to be used by all CLI commands
func authVGC() {
	// Initialize a new vault client
	vclient, err := vapi.NewClient(vapi.DefaultConfig())
	if err != nil {
		fmt.Println(errors.Wrap(err, "Failed to create vault client"))
		os.Exit(1)
	} else if vclient.Token() == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(errors.Wrap(err, "Could not find home directory to check ~/.vault-token"))
		} else {
			token, err := ioutil.ReadFile(home + "/.vault-token")
			if err != nil {
				if strings.Contains(err.Error(), "no such file or directory") {
					fmt.Println("INFO: Attempted to read token at ~/.vault-token, but the file does not exist")
					fmt.Println("Could not find token at VAULT_TOKEN or ~/.vault-token, exiting")
					os.Exit(1)
				} else {
					fmt.Println(errors.Wrap(err, "Failed to read ~/.vault-token file"))
					os.Exit(1)
				}
			}
			vclient.SetToken(string(token))
		}
	}

	// Add the Vault client to the Vaku client
	vgc = vaku.NewClient()
	vgc.Client = vclient
}

func print(i map[string]interface{}) {
	if format == "json" {
		json, err := json.MarshalIndent(i, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	} else if format == "text" {
		for _, v := range i {
			textPrint(v)
		}
	} else {
		fmt.Printf("ERROR: %s is not a valid or supported output format", format)
	}
}

func textPrint(i interface{}) {
	switch t := i.(type) {
	case map[string]map[string]interface{}:
		for k, v := range t {
			fmt.Printf("\n%+v\n", k)
			fmt.Println(strings.Repeat("-", len(k)))
			textPrint(v)
		}
	case map[string]interface{}:
		for k, v := range t {
			fmt.Printf("%s => %+v\n", k, v)
		}
	case []string:
		for _, v := range t {
			fmt.Println(v)
		}
	case string:
		fmt.Println(t)
	default:
		fmt.Printf("%+v\n", t)
	}
}

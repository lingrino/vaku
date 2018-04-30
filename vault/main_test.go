package vault

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	var err error

	V := NewClient()
	err = V.simpleInit()
	if err != nil {
		log.Fatalf("[FATAL]: TestMain: Failed to init the vault client: %s", err)
	}

	err = V.seed()
	if err != nil {
		log.Fatalf("[FATAL]: TestMain: Failed to seed vault: %s", err)
	}

	os.Exit(m.Run())
}

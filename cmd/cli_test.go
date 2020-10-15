package cmd

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCLI(t *testing.T) {
	t.Parallel()

	cli := newCLI()
	assert.Nil(t, cli.vc)
	assert.NotNil(t, cli.cmd)
	assert.Equal(t, "", cli.version)

	cli.setVersion("1.0.0")
	assert.Equal(t, "1.0.0", cli.version)
}

func TestInitVakuClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		giveFail    string
		giveDstAddr string
		giveWorkers int
		wantErr     string
	}{
		{
			name:        "valid",
			giveDstAddr: "foo:8200",
			giveWorkers: 1,
			wantErr:     "",
		},
		{
			name:        "bad addr",
			giveDstAddr: "\n",
			giveWorkers: 1,
			wantErr:     "initializing vaku client\nsetting vault address\nfailed to set address: parse \"\\n\": net/url: invalid control character in URL", //nolint:lll
		},
		{
			name:        "bad workers",
			giveWorkers: 0,
			wantErr:     "initializing vaku client\nworkers must 1 or greater: 0: invalid workers",
		},
		{
			name:        "fail vault.NewClient",
			giveFail:    "vault.NewClient",
			giveWorkers: 1,
			wantErr:     "initializing vaku client\ncreating new vault client",
		},
		{
			name:        "fail config.DefaultTokenHelper",
			giveFail:    "config.DefaultTokenHelper",
			giveWorkers: 1,
			wantErr:     "initializing vaku client\nsetting vault token\ngetting default token helper",
		},
		{
			name:        "fail helper.Get",
			giveFail:    "helper.Get",
			giveWorkers: 1,
			wantErr:     "initializing vaku client\nsetting vault token\nusing helper to get vault token",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cli, _, _ := newTestCLI(t, nil)
			cli.flagDstToken = "token"

			cli.fail = tt.giveFail
			cli.flagDstAddr = tt.giveDstAddr
			cli.flagWorkers = tt.giveWorkers

			err := cli.initVakuClient(cli.cmd, nil)

			errStr := ""
			if err != nil {
				errStr = err.Error()
			}

			assert.Equal(t, tt.wantErr, errStr)
		})
	}
}

func TestExecute(t *testing.T) {
	t.Parallel()

	var outW, errW bytes.Buffer

	code := Execute("dev", os.Args[1:], &outW, &errW)
	assert.Equal(t, exitSuccess, code)

	code = Execute("dev", []string{"INVALID"}, &outW, &errW)
	assert.Equal(t, exitFailure, code)
}

// TestHasExample tests that every command has an example.
func TestHasExample(t *testing.T) {
	t.Parallel()

	cli, _, _ := newTestCLI(t, nil)
	assert.True(t, allHasExample(cli.cmd))
}

// allHasExample recursively checks a command and it's children for example functions.
func allHasExample(cmds ...*cobra.Command) bool {
	res := true
	for _, cmd := range cmds {
		res = res && (cmd.HasExample() || cmd.Hidden) && allHasExample(cmd.Commands()...)
	}
	return res
}

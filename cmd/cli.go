package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/config"
	"github.com/lingrino/vaku/vaku"
	"github.com/spf13/cobra"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

var (
	errInitVakuClient   = errors.New("error initializing vaku client")
	errNewVaultClient   = errors.New("error creating new vault client")
	errVaultReadEnv     = errors.New("error with vault reading the environment")
	errVaultTokenHelper = errors.New("error getting default token helper")
	errGetVaultToken    = errors.New("error using helper to get vault token")
	errSetVaultToken    = errors.New("error setting vault token")
)

// cli extends cobra.Command with our own config.
type cli struct {
	// clients
	vc  *vaku.Client
	cmd *cobra.Command

	// flags
	flagAbsPath bool
	flagFormat  string
	flagIndent  string
	flagSort    bool
	flagWorkers int

	// vault flags
	flagSrcAddr  string
	flagSrcToken string
	flagDstAddr  string
	flagDstToken string

	// data
	version string
}

// newCLI returns a new CLI ready to run. Vaku client is not set because some commands (version) do
// not need it. Instead vc is initialized as a persistent function on the path/folder subcommands.
func newCLI() *cli {
	cli := &cli{}
	cli.cmd = cli.newVakuCmd()
	return cli
}

// setVersion sets the CLI version.
func (c *cli) setVersion(version string) {
	c.version = version
}

// initVakuClient initializes our vaku client and underlying vault clients.
// https://github.com/hashicorp/vault/blob/8571221f03c92ac3acac27c240fa7c9b3cb22db5/command/base.go#L67-L159
func (c *cli) initVakuClient(cmd *cobra.Command, args []string) error {
	var options []vaku.Option

	srcClient, err := c.newVaultClient(c.flagSrcAddr, c.flagSrcToken)
	if err != nil {
		return errInitVakuClient
	}
	options = append(options, vaku.WithVaultSrcClient(srcClient))

	if c.flagDstAddr != "" || c.flagDstToken != "" {
		dstClient, err := c.newVaultClient(c.flagDstAddr, c.flagDstToken)
		if err != nil {
			return errInitVakuClient
		}
		options = append(options, vaku.WithVaultDstClient(dstClient))
	}

	options = append(options, vaku.WithabsolutePath(c.flagAbsPath))
	options = append(options, vaku.WithWorkers(c.flagWorkers))

	vakuClient, err := vaku.NewClient(options...)
	if err != nil {
		return errInitVakuClient
	}

	c.vc = vakuClient

	return nil
}

// newVaultClient creates a new vault client. Prefer passed addr/token. Fallback to env/config.
func (c *cli) newVaultClient(addr, token string) (*api.Client, error) {
	cfg := api.DefaultConfig()
	err := cfg.ReadEnvironment()
	if err != nil {
		return nil, errVaultReadEnv
	}

	if addr != "" {
		cfg.Address = addr
	}

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, errNewVaultClient
	}

	err = c.setVaultToken(client, token)
	if err != nil {
		return nil, errSetVaultToken
	}

	if os.Getenv(api.EnvVaultMaxRetries) == "" {
		client.SetMaxRetries(0)
	}

	return client, nil
}

// setVaultToken sets vault token on client. Prefer passed token. Fallback to env/config.
func (c *cli) setVaultToken(vc *api.Client, token string) error {
	if token != "" {
		vc.SetToken(token)
		return nil
	}
	token = vc.Token()
	if token == "" {
		helper, err := config.DefaultTokenHelper()
		if err != nil {
			return errVaultTokenHelper
		}
		token, err = helper.Get()
		if err != nil {
			return errGetVaultToken
		}
		vc.SetToken(token)
	}
	return nil
}

// Execute runs the CLI.
func Execute(version string, args []string, outW, errW io.Writer) int {
	cli := newCLI()
	cli.setVersion(version)

	cli.cmd.SetArgs(args)
	cli.cmd.SetOut(outW)
	cli.cmd.SetErr(errW)
	err := cli.cmd.Execute()
	if err != nil {
		fmt.Fprintf(cli.cmd.ErrOrStderr(), "Error: %s\n", err)
		return exitFailure
	}

	return exitSuccess
}

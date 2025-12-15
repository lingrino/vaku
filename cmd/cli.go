package cmd

import (
	"errors"
	"io"
	"os"

	"github.com/hashicorp/vault/api/cliconfig"
	"github.com/spf13/cobra"

	vault "github.com/hashicorp/vault/api"
	vaku "github.com/lingrino/vaku/v2/api"
)

const (
	exitSuccess = 0
	exitFailure = 1
)

var (
	errInitVakuClient   = errors.New("initializing vaku client")
	errNewVaultClient   = errors.New("creating new vault client")
	errVaultTokenHelper = errors.New("getting default token helper")
	errGetVaultToken    = errors.New("using helper to get vault token")
	errSetVaultToken    = errors.New("setting vault token")
	errSetAddress       = errors.New("setting vault address")
)

// cli extends cobra.Command with our own config.
type cli struct {
	// clients
	vc  vaku.ClientInterface
	cmd *cobra.Command

	// flags
	flagAbsPath     bool
	flagNoAccessErr bool
	flagFormat      string
	flagIndent      string
	flagSort        bool
	flagWorkers     int

	// vault flags
	flagSrcAddr         string
	flagSrcToken        string
	flagSrcNspc         string
	flagDstAddr         string
	flagDstToken        string
	flagDstNspc         string
	flagMountPath       string
	flagMountVersion    string
	flagSrcMountPath    string
	flagSrcMountVersion string
	flagDstMountPath    string
	flagDstMountVersion string

	// data
	version string

	// failure injection
	fail string
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

// initVakuClient initializes a vaku client in the cli struct
// https://github.com/hashicorp/vault/blob/8571221f03c92ac3acac27c240fa7c9b3cb22db5/command/base.go#L67-L159
func (c *cli) initVakuClient(cmd *cobra.Command, args []string) error {
	// validate flags first (child PersistentPreRunE overrides parent's validateVakuFlags)
	if err := c.validateVakuFlags(cmd, args); err != nil {
		return err
	}

	// don't proceed if vc is already set (likely in tests)
	if c.vc != nil {
		return nil
	}

	vc, err := c.newVakuClient()
	if err != nil {
		return err
	}

	c.vc = vc

	return nil
}

// newVakuClient creates a vaku client and underlying vault clients.
func (c *cli) newVakuClient() (*vaku.Client, error) {
	var options []vaku.Option

	srcClient, err := c.newVaultClient(c.flagSrcAddr, c.flagSrcNspc, c.flagSrcToken)
	if err != nil {
		return nil, c.combineErr(errInitVakuClient, err)
	}
	options = append(options, vaku.WithVaultSrcClient(srcClient))

	if c.flagDstAddr != "" || c.flagDstToken != "" {
		dstClient, err := c.newVaultClient(c.flagDstAddr, c.flagDstNspc, c.flagDstToken)
		if err != nil {
			return nil, c.combineErr(errInitVakuClient, err)
		}
		options = append(options, vaku.WithVaultDstClient(dstClient))
	}

	options = append(options, vaku.WithAbsolutePath(c.flagAbsPath))
	options = append(options, vaku.WithIgnoreAccessErrors(c.flagNoAccessErr))
	options = append(options, vaku.WithWorkers(c.flagWorkers))

	// Source mount provider - use explicit source flags or short aliases
	srcMountPath := c.getSrcMountPath()
	if srcMountPath != "" {
		srcMountVersion := c.getSrcMountVersion()
		options = append(options,
			vaku.WithSrcMountProvider(vaku.NewStaticMountProvider(srcMountPath, srcMountVersion)))
	}

	// Destination mount provider
	if c.flagDstMountPath != "" {
		options = append(options,
			vaku.WithDstMountProvider(vaku.NewStaticMountProvider(c.flagDstMountPath, c.flagDstMountVersion)))
	}

	vakuClient, err := vaku.NewClient(options...)
	if err != nil {
		return nil, c.combineErr(errInitVakuClient, err)
	}

	return vakuClient, nil
}

// newVaultClient creates a new vault client. Prefer passed addr/token. Fallback to env/config.
func (c *cli) newVaultClient(addr, namespace, token string) (*vault.Client, error) {
	// nil means use default configuration and read from environment
	client, err := vault.NewClient(nil)
	if err != nil || c.fail == "vault.NewClient" {
		return nil, c.combineErr(errNewVaultClient, err)
	}

	if addr != "" {
		err := client.SetAddress(addr)
		if err != nil {
			return nil, c.combineErr(errSetAddress, err)
		}
	}

	if namespace != "" {
		client.SetNamespace(namespace)
	}

	err = c.setVaultToken(client, token)
	if err != nil {
		return nil, c.combineErr(errSetVaultToken, err)
	}

	if os.Getenv(vault.EnvVaultMaxRetries) == "" {
		client.SetMaxRetries(0)
	}

	return client, nil
}

// setVaultToken sets vault token on client. Prefer passed token. Fallback to env/config.
func (c *cli) setVaultToken(vc *vault.Client, token string) error {
	if token != "" {
		vc.SetToken(token)
	}
	if vc.Token() != "" {
		return nil
	}

	helper, err := cliconfig.DefaultTokenHelper()
	if err != nil || c.fail == "config.DefaultTokenHelper" {
		return c.combineErr(errVaultTokenHelper, err)
	}
	token, err = helper.Get()
	if err != nil || c.fail == "helper.Get" {
		return c.combineErr(errGetVaultToken, err)
	}
	vc.SetToken(token)

	return nil
}

// Execute runs a standard CLI and can be called externally.
func Execute(version string, args []string, outW, errW io.Writer) int {
	cli := newCLI()
	cli.setVersion(version)

	cli.cmd.SetArgs(args)
	cli.cmd.SetOut(outW)
	cli.cmd.SetErr(errW)

	return cli.execute()
}

// execute runs the CLI. Expects args and out/err writers to be set.
func (c *cli) execute() int {
	err := c.cmd.Execute()
	if err != nil {
		c.output(err)
		return exitFailure
	}

	return exitSuccess
}

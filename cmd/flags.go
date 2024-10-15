package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

// Base Flags.
const (
	flagFormatName    = "format"
	flagFormatUse     = "output format: text|json"
	flagFormatDefault = "text"

	flagIndentName    = "indent-char"
	flagIndentShort   = "i"
	flagIndentUse     = "string used for indents"
	flagIndentDefault = "    "

	flagSortName    = "sort"
	flagSortShort   = "s"
	flagSortUse     = "sort output text"
	flagSortDefault = true
)

// Vault Flags.
const (
	flagAbsPathName    = "absolute-path"
	flagAbsPathShort   = "p"
	flagAbsPathUse     = "show absolute path in output"
	flagAbsPathDefault = false

	flagWorkersName    = "workers"
	flagWorkersShort   = "w"
	flagWorkersUse     = "number of concurrent workers"
	flagWorkersDefault = 10

	flagAddrName    = "address"
	flagAddrShort   = "a"
	flagAddrUse     = "address of the Vault server"
	flagAddrDefault = ""

	flagSrcAddrName    = "source-address"
	flagSrcAddrUse     = "address of the source Vault server (alias for --address)"
	flagSrcAddrDefault = ""

	flagDstAddrName    = "destination-address"
	flagDstAddrUse     = "address of the destination Vault server"
	flagDstAddrDefault = ""

	flagNspcName    = "namespace"
	flagNspcShort   = "n"
	flagNspcUse     = "name of the vault namespace to use in the source client"
	flagNspcDefault = ""

	flagSrcNspcName    = "source-namespace"
	flagSrcNspcUse     = "name of the vault namespace to use in the source client (alias for --namespace)"
	flagSrcNspcDefault = ""

	flagDstNspcName    = "destination-namespace"
	flagDstNspcUse     = "name of the vault namespace to use in the destination client"
	flagDstNspcDefault = ""

	flagTokenName    = "token"
	flagTokenShort   = "t"
	flagTokenUse     = "token for the vault server"
	flagTokenDefault = ""

	flagSrcTokenName    = "source-token"
	flagSrcTokenUse     = "token for the source vault server (alias for --token)"
	flagSrcTokenDefault = ""

	flagDstTokenName    = "destination-token"
	flagDstTokenUse     = "token for the destination vault server (alias for --token)"
	flagDstTokenDefault = ""

	flagIgnoreErrorName    = "ignore-error"
	flagIgnoreErrorShort   = "q"
	flagIgnoreErrorUse     = "to ignore permission errors"
	flagIgnoreErrorDefault = false
)

var (
	errFlagInvalidFormat  = errors.New("format must be one of: text|json")
	errFlagInvalidWorkers = errors.New("workers must be >= 1")
)

// addVakuFlags adds all flags for the vaku command.
func (c *cli) addVakuFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&c.flagFormat, flagFormatName, flagFormatDefault, flagFormatUse)
	cmd.PersistentFlags().StringVarP(&c.flagIndent, flagIndentName, flagIndentShort, flagIndentDefault, flagIndentUse)
	cmd.PersistentFlags().BoolVarP(&c.flagSort, flagSortName, flagSortShort, flagSortDefault, flagSortUse)
}

// addPathFolderFlags adds all flags for the path and folder commands.
func (c *cli) addPathFolderFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&c.flagAbsPath, flagAbsPathName, flagAbsPathShort, flagAbsPathDefault, flagAbsPathUse)
	cmd.PersistentFlags().IntVarP(&c.flagWorkers, flagWorkersName, flagWorkersShort, flagWorkersDefault, flagWorkersUse)

	cmd.PersistentFlags().StringVarP(&c.flagSrcAddr, flagAddrName, flagAddrShort, flagAddrDefault, flagAddrUse)
	cmd.PersistentFlags().StringVar(&c.flagSrcAddr, flagSrcAddrName, flagSrcAddrDefault, flagSrcAddrUse)
	cmd.PersistentFlags().StringVar(&c.flagDstAddr, flagDstAddrName, flagDstAddrDefault, flagDstAddrUse)

	cmd.PersistentFlags().StringVarP(&c.flagSrcNspc, flagNspcName, flagNspcShort, flagNspcDefault, flagNspcUse)
	cmd.PersistentFlags().StringVar(&c.flagSrcNspc, flagSrcNspcName, flagSrcNspcDefault, flagSrcNspcUse)
	cmd.PersistentFlags().StringVar(&c.flagDstNspc, flagDstNspcName, flagDstNspcDefault, flagDstNspcUse)

	cmd.PersistentFlags().StringVarP(&c.flagSrcToken, flagTokenName, flagTokenShort, flagTokenDefault, flagTokenUse)
	cmd.PersistentFlags().StringVar(&c.flagSrcToken, flagSrcTokenName, flagSrcTokenDefault, flagSrcTokenUse)
	cmd.PersistentFlags().StringVar(&c.flagDstToken, flagDstTokenName, flagDstTokenDefault, flagDstTokenUse)
	cmd.PersistentFlags().BoolVarP(&c.flagIgnoreError, flagIgnoreErrorName, flagIgnoreErrorShort, flagIgnoreErrorDefault, flagIgnoreErrorUse)
}

// validateFlags checks if valid flag values were passed. Use as cmd.PersistentPreRunE.
func (c *cli) validateVakuFlags(cmd *cobra.Command, args []string) error {
	validationFuncs := []func() error{
		c.validFormat,
		c.validWorkers,
	}

	for _, f := range validationFuncs {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

// validFormat checks if the format flag is a valid option.
func (c *cli) validFormat() error {
	validFormats := []string{"text", "json"}

	for _, v := range validFormats {
		if c.flagFormat == v {
			return nil
		}
	}
	return errFlagInvalidFormat
}

// validWorkers checks if the workers flag is a valid option.
func (c *cli) validWorkers() error {
	if c.flagWorkers < 1 {
		return errFlagInvalidWorkers
	}
	return nil
}

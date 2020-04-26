package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

const (
	flagAbsPathName    = "absolute-path"
	flagAbsPathShort   = "a"
	flagAbsPathUse     = "show absolute path in output"
	flagAbsPathDefault = false

	flagFormatName    = "format"
	flagFormatUse     = "output format: text|json"
	flagFormatDefault = "text"

	flagWorkersName    = "workers"
	flagWorkersShort   = "w"
	flagWorkersUse     = "number of concurrent workers"
	flagWorkersDefault = 10
)

var (
	errFlagInvalidFormat  = errors.New("format must be one of: text|json")
	errFlagInvalidWorkers = errors.New("workers must be >= 1")
)

// addVakuFlags adds all flags to the vaku command.
func (c *cli) addVakuFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVarP(&c.flagAbsPath, flagAbsPathName, flagAbsPathShort, flagAbsPathDefault, flagAbsPathUse)
	cmd.PersistentFlags().StringVar(&c.flagFormat, flagFormatName, flagFormatDefault, flagFormatUse)
	cmd.PersistentFlags().IntVarP(&c.flagWorkers, flagWorkersName, flagWorkersShort, flagWorkersDefault, flagWorkersUse)
}

// validateFlags checks if valid flag values were passed. Use as cmd.PersistentPreRunE
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

// validFormat checks if the format flag is a valid option
func (c *cli) validFormat() error {
	validFormats := []string{"text", "json"}

	for _, v := range validFormats {
		if c.flagFormat == v {
			return nil
		}
	}
	return errFlagInvalidFormat
}

// validWorkers checks if the workers flag is a valid option
func (c *cli) validWorkers() error {
	if c.flagWorkers < 1 {
		return errFlagInvalidWorkers
	}
	return nil
}

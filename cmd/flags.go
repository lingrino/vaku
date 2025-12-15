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

	flagNoAccessErrName    = "ignore-read-errors"
	flagNoAccessErrUse     = "ignore path read errors and continue"
	flagNoAccessErrDefault = false

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

	flagMountPathName    = "mount-path"
	flagMountPathShort   = "m"
	flagMountPathUse     = "source mount path (bypasses sys/mounts lookup, alias for --mount-path-source)"
	flagMountPathDefault = ""

	flagMountVersionName    = "mount-version"
	flagMountVersionUse     = "source mount version: 1|2 (requires --mount-path, alias for --mount-version-source)"
	flagMountVersionDefault = "2"

	flagSrcMountPathName    = "mount-path-source"
	flagSrcMountPathUse     = "source mount path (bypasses sys/mounts lookup)"
	flagSrcMountPathDefault = ""

	flagSrcMountVersionName    = "mount-version-source"
	flagSrcMountVersionUse     = "source mount version: 1|2 (requires --mount-path-source)"
	flagSrcMountVersionDefault = "2"

	flagDstMountPathName    = "mount-path-destination"
	flagDstMountPathUse     = "destination mount path (bypasses sys/mounts lookup)"
	flagDstMountPathDefault = ""

	flagDstMountVersionName    = "mount-version-destination"
	flagDstMountVersionUse     = "destination mount version: 1|2 (requires --mount-path-destination)"
	flagDstMountVersionDefault = "2"
)

var (
	errFlagInvalidFormat          = errors.New("format must be one of: text|json")
	errFlagInvalidWorkers         = errors.New("workers must be >= 1")
	errFlagInvalidMountVersion    = errors.New("mount-version must be one of: 1|2")
	errFlagMountVersionNoPath     = errors.New("mount-version requires --mount-path")
	errFlagInvalidSrcMountVersion = errors.New("mount-version-source must be one of: 1|2")
	errFlagSrcMountVersionNoPath  = errors.New("mount-version-source requires --mount-path-source")
	errFlagInvalidDstMountVersion = errors.New("mount-version-destination must be one of: 1|2")
	errFlagDstMountVersionNoPath  = errors.New("mount-version-destination requires --mount-path-destination")
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
	cmd.PersistentFlags().BoolVar(&c.flagNoAccessErr, flagNoAccessErrName, flagNoAccessErrDefault, flagNoAccessErrUse)
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

	cmd.PersistentFlags().StringVarP(&c.flagMountPath,
		flagMountPathName, flagMountPathShort, flagMountPathDefault, flagMountPathUse)
	cmd.PersistentFlags().StringVar(&c.flagMountVersion,
		flagMountVersionName, flagMountVersionDefault, flagMountVersionUse)

	cmd.PersistentFlags().StringVar(&c.flagSrcMountPath,
		flagSrcMountPathName, flagSrcMountPathDefault, flagSrcMountPathUse)
	cmd.PersistentFlags().StringVar(&c.flagSrcMountVersion,
		flagSrcMountVersionName, flagSrcMountVersionDefault, flagSrcMountVersionUse)

	cmd.PersistentFlags().StringVar(&c.flagDstMountPath,
		flagDstMountPathName, flagDstMountPathDefault, flagDstMountPathUse)
	cmd.PersistentFlags().StringVar(&c.flagDstMountVersion,
		flagDstMountVersionName, flagDstMountVersionDefault, flagDstMountVersionUse)
}

// validateFlags checks if valid flag values were passed. Use as cmd.PersistentPreRunE.
func (c *cli) validateVakuFlags(cmd *cobra.Command, args []string) error {
	validationFuncs := []func() error{
		c.validFormat,
		c.validWorkers,
		c.validMountFlags,
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

// isValidMountVersion checks if the version string is valid (1 or 2).
func isValidMountVersion(version string) bool {
	return version == "1" || version == "2"
}

// validateMountPair validates a mount path and version pair.
func validateMountPair(path, version, defaultVersion string, errNoPath, errInvalidVersion error) error {
	// If version is explicitly set to non-default and non-empty but path is empty, error
	if version != defaultVersion && version != "" && path == "" {
		return errNoPath
	}
	// If path is set, validate version
	if path != "" && !isValidMountVersion(version) {
		return errInvalidVersion
	}
	return nil
}

// validMountFlags checks if the mount flags are valid.
func (c *cli) validMountFlags() error {
	// Validate --mount-path / --mount-version (short aliases)
	if err := validateMountPair(c.flagMountPath, c.flagMountVersion, flagMountVersionDefault,
		errFlagMountVersionNoPath, errFlagInvalidMountVersion); err != nil {
		return err
	}

	// Validate --mount-path-source / --mount-version-source
	if err := validateMountPair(c.flagSrcMountPath, c.flagSrcMountVersion, flagSrcMountVersionDefault,
		errFlagSrcMountVersionNoPath, errFlagInvalidSrcMountVersion); err != nil {
		return err
	}

	// Validate --mount-path-destination / --mount-version-destination
	if err := validateMountPair(c.flagDstMountPath, c.flagDstMountVersion, flagDstMountVersionDefault,
		errFlagDstMountVersionNoPath, errFlagInvalidDstMountVersion); err != nil {
		return err
	}

	return nil
}

// getSrcMountPath returns the source mount path, preferring the explicit source flag over the alias.
func (c *cli) getSrcMountPath() string {
	if c.flagSrcMountPath != "" {
		return c.flagSrcMountPath
	}
	return c.flagMountPath
}

// getSrcMountVersion returns the source mount version, preferring the explicit source flag over alias.
func (c *cli) getSrcMountVersion() string {
	if c.flagSrcMountPath != "" {
		return c.flagSrcMountVersion
	}
	return c.flagMountVersion
}

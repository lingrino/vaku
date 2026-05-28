//! Clap argument tree for `vaku`. Flag names and short forms mirror the Go
//! cobra setup so existing scripts keep working unchanged.

use clap::{ArgAction, Args, Parser, Subcommand, ValueEnum};
use std::path::PathBuf;

pub const VAKU_LONG: &str = "Vaku is a CLI for working with large Vault k/v secret engines\n\n\
The Vaku CLI provides path- and folder-based commands that work on\n\
both Version 1 and Version 2 K/V secret engines. Vaku can help manage\n\
large amounts of Vault data by updating secrets in place, moving\n\
paths or folders, searching secrets, and more.\n\n\
Vaku is not a replacement for the Vault CLI and requires that you are\n\
already authenticated to Vault before running any commands. Vaku\n\
commands should not be run on non-K/V engines.\n\n\
CLI documentation - 'vaku help [cmd]'\n\
API documentation - https://docs.rs/vaku\n\
Built by Sean Lingren <sean@lingren.com>";

pub const VAKU_SHORT: &str = "Vaku is a CLI for working with large Vault k/v secret engines";
pub const VAKU_EXAMPLE: &str = "vaku folder list secret/foo";

pub const PATH_SHORT: &str = "Commands that act on Vault paths";
pub const PATH_LONG: &str = "Commands that act on Vault paths\n\n\
Commands under the path subcommand act on Vault paths. Vaku can list,\n\
copy, move, search, etc.. on Vault paths.";
pub const PATH_EXAMPLE: &str = "vaku path list secret/foo";

pub const FOLDER_SHORT: &str = "Commands that act on Vault folders";
pub const FOLDER_LONG: &str = "Commands that act on Vault folders\n\n\
Commands under the folder subcommand act on Vault folders. Folders\n\
are designated by paths that end in a '/' such as 'secret/foo/'. Vaku\n\
can list, copy, move, search, etc.. on Vault folders.";
pub const FOLDER_EXAMPLE: &str = "vaku folder list secret/foo";

pub const VERSION_SHORT: &str = "Print vaku version";
pub const VERSION_EXAMPLE: &str = "vaku version";

#[derive(Parser, Debug, Clone)]
#[command(
    name = "vaku",
    bin_name = "vaku",
    about = VAKU_SHORT,
    long_about = VAKU_LONG,
    after_help = VAKU_EXAMPLE,
    disable_help_subcommand = true,
    disable_version_flag = true,
    arg_required_else_help = false,
)]
pub struct VakuArgs {
    /// output format: text|json
    #[arg(
        long = "format",
        global = true,
        default_value = "text",
        help_heading = "Options"
    )]
    pub format: String,

    /// string used for indents
    #[arg(
        short = 'i',
        long = "indent-char",
        global = true,
        default_value = "    ",
        help_heading = "Options"
    )]
    pub indent_char: String,

    /// sort output text
    #[arg(short = 's', long = "sort", global = true, default_value_t = true, action = ArgAction::Set, help_heading = "Options")]
    pub sort: bool,

    #[command(subcommand)]
    pub cmd: Option<TopCmd>,
}

#[derive(Subcommand, Debug, Clone)]
pub enum TopCmd {
    /// Commands that act on Vault paths
    #[command(
        about = PATH_SHORT,
        long_about = PATH_LONG,
        after_help = PATH_EXAMPLE,
        subcommand_required = true,
        arg_required_else_help = true,
    )]
    Path(PathRoot),

    /// Commands that act on Vault folders
    #[command(
        about = FOLDER_SHORT,
        long_about = FOLDER_LONG,
        after_help = FOLDER_EXAMPLE,
        subcommand_required = true,
        arg_required_else_help = true,
    )]
    Folder(FolderRoot),

    /// Print vaku version
    #[command(about = VERSION_SHORT, after_help = VERSION_EXAMPLE)]
    Version,

    /// Generate markdown docs at a path
    #[command(hide = true, after_help = "vaku docs .")]
    Docs { path: PathBuf },

    /// Generate completion scripts (clap-built-in)
    #[command(hide = true)]
    Completion { shell: clap_complete::Shell },
}

#[derive(Args, Debug, Clone, Default)]
#[group(skip)]
pub struct PathFolderFlags {
    /// show absolute path in output
    #[arg(
        global = true,
        short = 'p',
        long = "absolute-path",
        help_heading = "Vault"
    )]
    pub absolute_path: bool,

    /// ignore path read errors and continue
    #[arg(global = true, long = "ignore-read-errors", help_heading = "Vault")]
    pub ignore_read_errors: bool,

    /// number of concurrent workers
    #[arg(
        global = true,
        short = 'w',
        long = "workers",
        default_value_t = 10,
        help_heading = "Vault"
    )]
    pub workers: usize,

    /// address of the Vault server
    #[arg(
        global = true,
        short = 'a',
        long = "address",
        default_value = "",
        help_heading = "Vault"
    )]
    pub address: String,

    /// address of the source Vault server (alias for --address)
    #[arg(
        global = true,
        long = "source-address",
        default_value = "",
        help_heading = "Vault"
    )]
    pub source_address: String,

    /// address of the destination Vault server
    #[arg(
        global = true,
        long = "destination-address",
        default_value = "",
        help_heading = "Vault"
    )]
    pub destination_address: String,

    /// name of the vault namespace to use in the source client
    #[arg(
        global = true,
        short = 'n',
        long = "namespace",
        default_value = "",
        help_heading = "Vault"
    )]
    pub namespace: String,

    /// name of the vault namespace to use in the source client (alias for --namespace)
    #[arg(
        global = true,
        long = "source-namespace",
        default_value = "",
        help_heading = "Vault"
    )]
    pub source_namespace: String,

    /// name of the vault namespace to use in the destination client
    #[arg(
        global = true,
        long = "destination-namespace",
        default_value = "",
        help_heading = "Vault"
    )]
    pub destination_namespace: String,

    /// token for the vault server
    #[arg(
        global = true,
        short = 't',
        long = "token",
        default_value = "",
        help_heading = "Vault"
    )]
    pub token: String,

    /// token for the source vault server (alias for --token)
    #[arg(
        global = true,
        long = "source-token",
        default_value = "",
        help_heading = "Vault"
    )]
    pub source_token: String,

    /// token for the destination vault server (alias for --token)
    #[arg(
        global = true,
        long = "destination-token",
        default_value = "",
        help_heading = "Vault"
    )]
    pub destination_token: String,

    /// source mount path (bypasses sys/mounts lookup, alias for --mount-path-source)
    #[arg(
        global = true,
        short = 'm',
        long = "mount-path",
        default_value = "",
        help_heading = "Vault"
    )]
    pub mount_path: String,

    /// source mount version: 1|2 (requires --mount-path, alias for --mount-version-source)
    #[arg(
        global = true,
        long = "mount-version",
        default_value = "2",
        help_heading = "Vault"
    )]
    pub mount_version: String,

    /// source mount path (bypasses sys/mounts lookup)
    #[arg(
        global = true,
        long = "mount-path-source",
        default_value = "",
        help_heading = "Vault"
    )]
    pub mount_path_source: String,

    /// source mount version: 1|2 (requires --mount-path-source)
    #[arg(
        global = true,
        long = "mount-version-source",
        default_value = "2",
        help_heading = "Vault"
    )]
    pub mount_version_source: String,

    /// destination mount path (bypasses sys/mounts lookup)
    #[arg(
        global = true,
        long = "mount-path-destination",
        default_value = "",
        help_heading = "Vault"
    )]
    pub mount_path_destination: String,

    /// destination mount version: 1|2 (requires --mount-path-destination)
    #[arg(
        global = true,
        long = "mount-version-destination",
        default_value = "2",
        help_heading = "Vault"
    )]
    pub mount_version_destination: String,
}

#[derive(Args, Debug, Clone)]
pub struct PathRoot {
    #[command(flatten)]
    pub flags: PathFolderFlags,
    #[command(subcommand)]
    pub cmd: PathSub,
}

#[derive(Args, Debug, Clone)]
pub struct FolderRoot {
    #[command(flatten)]
    pub flags: PathFolderFlags,
    #[command(subcommand)]
    pub cmd: FolderSub,
}

#[derive(Subcommand, Debug, Clone)]
pub enum PathSub {
    /// List all paths at a path
    #[command(after_help = "vaku path list secret/foo")]
    List { path: String },
    /// Read a secret at a path
    #[command(after_help = "vaku path read secret/foo")]
    Read { path: String },
    /// Delete a secret at a path
    #[command(after_help = "vaku path delete secret/foo")]
    Delete { path: String },
    /// Delete all secret metadata and versions at a path. V2 engines only.
    #[command(name = "delete-meta", after_help = "vaku path delete-meta secret/foo")]
    DeleteMeta { path: String },
    /// Search a secret for a search string
    #[command(after_help = "vaku path search secret/foo bar")]
    Search { path: String, search: String },
    /// Copy a secret from a source path to a destination path
    #[command(after_help = "vaku path copy secret/foo secret/bar")]
    Copy {
        src: String,
        dst: String,
        #[arg(
            long = "all-versions",
            help = "copy all versions of the secret (KV v2 only)"
        )]
        all_versions: bool,
    },
    /// Move a secret from a source path to a destination path
    #[command(after_help = "vaku path move secret/foo secret/bar")]
    Move {
        src: String,
        dst: String,
        #[arg(
            long = "all-versions",
            help = "move all versions of the secret (KV v2 only)"
        )]
        all_versions: bool,
        #[arg(
            long = "destroy",
            help = "permanently destroy all versions at source after copy (KV v2 only)"
        )]
        destroy: bool,
    },
    /// (hidden) Use the vaku API or native Vault CLI
    #[command(hide = true, disable_help_subcommand = true)]
    Write {
        #[arg(allow_hyphen_values = true)]
        args: Vec<String>,
    },
    /// (hidden) Use the vaku API or native Vault CLI
    #[command(hide = true, disable_help_subcommand = true)]
    Update {
        #[arg(allow_hyphen_values = true)]
        args: Vec<String>,
    },
    /// (hidden) Use the vaku API or native Vault CLI
    #[command(hide = true, disable_help_subcommand = true)]
    Destroy {
        #[arg(allow_hyphen_values = true)]
        args: Vec<String>,
    },
}

#[derive(Subcommand, Debug, Clone)]
pub enum FolderSub {
    /// Recursively list all paths in a folder
    #[command(after_help = "vaku folder list secret/foo")]
    List { path: String },
    /// Recursively read all secrets in a folder
    #[command(after_help = "vaku folder read secret/foo")]
    Read { path: String },
    /// Write a folder of secrets. WARNING: command expects a very specific json input
    #[command(after_help = "vaku folder write '{\"a/b/c\": {\"foo\": \"bar\"}}'")]
    Write { json: String },
    /// Recursively delete all secrets in a folder
    #[command(after_help = "vaku folder delete secret/foo")]
    Delete { path: String },
    /// Recursively delete all secrets metadata and versions in a folder. V2 engines only.
    #[command(
        name = "delete-meta",
        after_help = "vaku folder delete-meta secret/foo"
    )]
    DeleteMeta { path: String },
    /// Recursively search all secrets in a folder for a search string
    #[command(after_help = "vaku folder search secret/foo bar")]
    Search { path: String, search: String },
    /// Recursively copy all secrets in source folder to destination folder
    #[command(after_help = "vaku folder copy secret/foo secret/bar")]
    Copy {
        src: String,
        dst: String,
        #[arg(
            long = "all-versions",
            help = "copy all versions of the secret (KV v2 only)"
        )]
        all_versions: bool,
    },
    /// Recursively move all secrets in source folder to destination folder
    #[command(after_help = "vaku folder move secret/foo secret/bar")]
    Move {
        src: String,
        dst: String,
        #[arg(
            long = "all-versions",
            help = "move all versions of the secret (KV v2 only)"
        )]
        all_versions: bool,
        #[arg(
            long = "destroy",
            help = "permanently destroy all versions at source after copy (KV v2 only)"
        )]
        destroy: bool,
    },
    /// (hidden)
    #[command(hide = true, disable_help_subcommand = true)]
    Destroy {
        #[arg(allow_hyphen_values = true)]
        args: Vec<String>,
    },
}

/// Available output formats.
#[derive(Copy, Clone, Debug, PartialEq, Eq, ValueEnum)]
pub enum Format {
    Text,
    Json,
}

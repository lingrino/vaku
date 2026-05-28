//! Error type for the Vaku library.
//!
//! Mirrors the Go `wrapErr` semantics: every error has a "kind" sentinel
//! (similar to a Go sentinel `error` value compared with `errors.Is`) plus an
//! optional wrapped source error. Walk the chain with [`Error::source`] and
//! the public [`compare_errors`] helper.

use std::error::Error as StdError;
use std::fmt;

/// Boxed-trait-object error suitable for wrapping arbitrary causes.
pub type BoxError = Box<dyn StdError + Send + Sync + 'static>;

/// Result alias used throughout the crate.
pub type Result<T> = std::result::Result<T, Error>;

/// Kinds of errors Vaku can produce. Each variant maps 1:1 to a sentinel
/// `ErrXxx` in the Go implementation and is matched by tests via
/// [`compare_errors`].
#[derive(Debug, Clone, PartialEq, Eq, Hash)]
pub enum ErrorKind {
    // Generic
    Context,
    DecodeSecret,
    JsonMarshal,
    NilData,
    UnknownError,
    NumWorkers,
    ApplyOptions,
    // Mount
    MountInfo,
    ListMounts,
    NoMount,
    RewritePath,
    MountVersion,
    // Path
    PathList,
    VaultList,
    PathRead,
    PathReadVersion,
    VaultRead,
    PathReadMeta,
    PathWrite,
    VaultWrite,
    PathDelete,
    VaultDelete,
    PathDeleteMeta,
    PathDestroy,
    PathUpdate,
    PathSearch,
    PathCopy,
    PathCopyAllVersions,
    PathMove,
    PathMoveAllVersions,
    // Folder
    FolderList,
    FolderListChan,
    FolderRead,
    FolderReadChan,
    FolderWrite,
    FolderDelete,
    FolderDeleteMeta,
    FolderDestroy,
    FolderSearch,
    FolderCopy,
    FolderCopyAllVersions,
    FolderMove,
    FolderMoveAllVersions,
    /// "Custom" carries a free-form message (mirrors Go's `errors.New(msg)`
    /// inside `newWrapErr` when only `msg` was supplied).
    Custom(String),
}

impl ErrorKind {
    /// The Go-style sentinel message for this kind. The CLI and tests rely on
    /// these exact strings (e.g. `"path list"`).
    pub fn message(&self) -> &str {
        match self {
            ErrorKind::Context => "context",
            ErrorKind::DecodeSecret => "decode secret",
            ErrorKind::JsonMarshal => "json marshal",
            ErrorKind::NilData => "nil data",
            ErrorKind::UnknownError => "unknown error",
            ErrorKind::NumWorkers => "invalid workers",
            ErrorKind::ApplyOptions => "applying options",
            ErrorKind::MountInfo => "mount info",
            ErrorKind::ListMounts => "list mounts",
            ErrorKind::NoMount => "no matching mount",
            ErrorKind::RewritePath => "rewriting path",
            ErrorKind::MountVersion => "mount version does not support operation",
            ErrorKind::PathList => "path list",
            ErrorKind::VaultList => "vault list",
            ErrorKind::PathRead => "path read",
            ErrorKind::PathReadVersion => "path read version",
            ErrorKind::VaultRead => "vault read",
            ErrorKind::PathReadMeta => "path read meta",
            ErrorKind::PathWrite => "path write",
            ErrorKind::VaultWrite => "vault write",
            ErrorKind::PathDelete => "path delete",
            ErrorKind::VaultDelete => "vault delete",
            ErrorKind::PathDeleteMeta => "path delete meta",
            ErrorKind::PathDestroy => "path destroy",
            ErrorKind::PathUpdate => "path update",
            ErrorKind::PathSearch => "path search",
            ErrorKind::PathCopy => "path copy",
            ErrorKind::PathCopyAllVersions => "path copy all versions",
            ErrorKind::PathMove => "path move",
            ErrorKind::PathMoveAllVersions => "path move all versions",
            ErrorKind::FolderList => "folder list",
            ErrorKind::FolderListChan => "folder list chan",
            ErrorKind::FolderRead => "folder read",
            ErrorKind::FolderReadChan => "folder read chan",
            ErrorKind::FolderWrite => "folder write",
            ErrorKind::FolderDelete => "folder delete",
            ErrorKind::FolderDeleteMeta => "folder delete meta",
            ErrorKind::FolderDestroy => "folder destroy",
            ErrorKind::FolderSearch => "folder search",
            ErrorKind::FolderCopy => "folder copy",
            ErrorKind::FolderCopyAllVersions => "folder copy all versions",
            ErrorKind::FolderMove => "folder move",
            ErrorKind::FolderMoveAllVersions => "folder move all versions",
            ErrorKind::Custom(s) => s.as_str(),
        }
    }
}

/// A wrapping error mirroring Go's `wrapErr` struct.
pub struct Error {
    kind: ErrorKind,
    msg: String,
    source: Option<BoxError>,
}

impl Error {
    /// Construct a new error. Reproduces the message-fallback rules from
    /// Go's `newWrapErr`:
    ///
    /// | msg | kind | source | resulting kind          | resulting msg                           |
    /// |-----|------|--------|-------------------------|-----------------------------------------|
    /// | ""  | None | None   | UnknownError            | "unknown error"                         |
    /// | ""  | None | Some   | UnknownError            | "unknown error: <src>"                  |
    /// | "x" | None | _      | Custom("x")             | "x" or "x: <src>" if source             |
    /// | _   | Some | _      | <kind>                  | "<msg>: <kind.msg>[: <src>]" or shorter |
    pub fn new(msg: Option<String>, kind: Option<ErrorKind>, source: Option<BoxError>) -> Self {
        let (kind, msg_or_empty) = match (msg, kind) {
            (Some(m), None) if m.is_empty() => (ErrorKind::UnknownError, String::new()),
            (None, None) => (ErrorKind::UnknownError, String::new()),
            (Some(m), None) => (ErrorKind::Custom(m.clone()), m),
            (None, Some(k)) => (k, String::new()),
            (Some(m), Some(k)) => (k, m),
        };

        let kind_msg = kind.message();
        let msg: String = match (msg_or_empty.is_empty(), source.as_ref()) {
            (true, None) => kind_msg.to_string(),
            (true, Some(src)) => {
                if matches!(kind, ErrorKind::Custom(_)) {
                    format!("{kind_msg}: {src}")
                } else {
                    format!("{kind_msg}: {src}")
                }
            }
            (false, None) => {
                if msg_or_empty == kind_msg {
                    kind_msg.to_string()
                } else {
                    format!("{msg_or_empty}: {kind_msg}")
                }
            }
            (false, Some(src)) => {
                if msg_or_empty == kind_msg {
                    format!("{kind_msg}: {src}")
                } else {
                    format!("{msg_or_empty}: {kind_msg}: {src}")
                }
            }
        };

        Self { kind, msg, source }
    }

    /// Shortcut: wrap with kind only (and optional source).
    pub fn wrap(msg: &str, kind: ErrorKind, source: Option<BoxError>) -> Self {
        let m = if msg.is_empty() { None } else { Some(msg.to_string()) };
        Self::new(m, Some(kind), source)
    }

    /// Shortcut: wrap arbitrary text (Go's `errors.New`-style).
    pub fn from_msg(msg: impl Into<String>) -> Self {
        Self::new(Some(msg.into()), None, None)
    }

    /// Wrap an underlying error in [`ErrorKind::Context`] (mirrors Go's
    /// `ctxErr`). Returns `None` when `err` is `None`.
    pub(crate) fn ctx(err: Option<BoxError>) -> Option<Self> {
        err.map(|e| Self::new(None, Some(ErrorKind::Context), Some(e)))
    }

    /// The error kind (the topmost sentinel).
    pub fn kind(&self) -> &ErrorKind {
        &self.kind
    }
}

impl fmt::Debug for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.debug_struct("Error")
            .field("kind", &self.kind)
            .field("msg", &self.msg)
            .field("source", &self.source.as_ref().map(|s| s.to_string()))
            .finish()
    }
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(&self.msg)
    }
}

impl StdError for Error {
    fn source(&self) -> Option<&(dyn StdError + 'static)> {
        self.source.as_deref().map(|s| s as _)
    }
}

/// Matches one node of an error chain. A [`Plain`] matcher matches by kind
/// (equality); a [`Custom`] matcher matches by free-form message.
#[derive(Debug, Clone)]
pub enum ErrMatch {
    Plain(ErrorKind),
    Custom(String),
}

impl From<ErrorKind> for ErrMatch {
    fn from(k: ErrorKind) -> Self {
        ErrMatch::Plain(k)
    }
}

/// Walk an error chain confirming each node matches the corresponding matcher,
/// then assert the chain has terminated. Mirrors the Go test helper
/// `compareErrors` exactly.
///
/// Panics with a diagnostic message on the first mismatch — intended for use
/// from tests.
pub fn compare_errors(err: Option<&(dyn StdError + 'static)>, expected: &[ErrMatch]) {
    let mut current: Option<&(dyn StdError + 'static)> = err;
    for (i, want) in expected.iter().enumerate() {
        let Some(node) = current else {
            panic!("error chain ended before matcher {i} ({want:?})");
        };

        let matches = match (want, node.downcast_ref::<Error>()) {
            (ErrMatch::Plain(want_kind), Some(e)) => matches_kind(e.kind(), want_kind),
            (ErrMatch::Custom(want_msg), Some(e)) => match e.kind() {
                ErrorKind::Custom(have) => have == want_msg,
                _ => false,
            },
            // Foreign error in chain: compare via Display.
            (ErrMatch::Custom(want_msg), None) => node.to_string() == *want_msg,
            (ErrMatch::Plain(_), None) => false,
        };
        if !matches {
            panic!("error chain mismatch at depth {i}: wanted {want:?}, got {node}");
        }

        current = node.source();
    }
    if let Some(extra) = current {
        panic!("expected error chain to end but found extra node: {extra}");
    }
}

fn matches_kind(have: &ErrorKind, want: &ErrorKind) -> bool {
    match (have, want) {
        (ErrorKind::Custom(a), ErrorKind::Custom(b)) => a == b,
        (a, b) => std::mem::discriminant(a) == std::mem::discriminant(b),
    }
}

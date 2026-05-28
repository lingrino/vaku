//! Vaku is a library that extends HashiCorp Vault's K/V secret engine with
//! recursive (folder-level) and convenience (path-level) operations.
//!
//! The public API mirrors the original Go library 1:1 in semantics.

pub mod api;
pub mod cli;

pub use api::client::{Client, ClientBuilder};
pub use api::error::{Error, ErrorKind, Result};
pub use api::helpers::*;
pub use api::logical::{Logical, Secret, VaultHttpClient};
pub use api::mount_provider::{Mount, MountProvider, StaticMountProvider};
pub use api::secret::{SecretMeta, SecretVersionMeta};
pub use api::version::version;

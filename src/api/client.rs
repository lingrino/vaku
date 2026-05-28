//! The Vaku [`Client`] and its builder.
//!
//! Mirrors the Go `Client` + functional-options API: the builder methods are
//! the Rust analogue of `WithXxx`. A client always has a source side and an
//! optional destination side (for copy/move operations). If no destination
//! client is supplied the same source client is used.

use crate::api::error::{Error, ErrorKind};
use crate::api::helpers::{add_prefix_list, ensure_prefix_map, path_join, trim_prefix_map};
use crate::api::logical::{Logical, VaultHttpClient};
use crate::api::mount_provider::{DefaultMountProvider, MountProvider};
use serde_json::{Map, Value};
use std::collections::BTreeMap;
use std::sync::Arc;

/// Default worker count for concurrent folder operations.
pub const DEFAULT_WORKERS: usize = 10;

/// The Vaku client. Cheap to clone (`Arc` internally).
#[derive(Clone)]
pub struct Client {
    pub(crate) inner: Arc<ClientInner>,
}

/// Inner shared client state. Boxed up so [`Client`] is `Clone`-able cheaply
/// and so workers can capture an [`Arc`] of the state.
pub(crate) struct ClientInner {
    pub(crate) src: Side,
    pub(crate) dst: Side,
    pub(crate) workers: usize,
    pub(crate) absolute_path: bool,
    pub(crate) ignore_access_errors: bool,
}

/// One side of a client — either source or destination.
pub struct Side {
    pub logical: Arc<dyn Logical>,
    pub mount_provider: Arc<dyn MountProvider>,
}

impl Client {
    /// Returns a fresh [`ClientBuilder`].
    pub fn builder() -> ClientBuilder {
        ClientBuilder::default()
    }

    pub fn src(&self) -> &Side {
        &self.inner.src
    }
    pub fn dst(&self) -> &Side {
        &self.inner.dst
    }
    pub fn workers(&self) -> usize {
        self.inner.workers
    }
    pub fn absolute_path(&self) -> bool {
        self.inner.absolute_path
    }
    pub fn ignore_access_errors(&self) -> bool {
        self.inner.ignore_access_errors
    }

    /// Returns a [`Client`] view that uses the destination side as its source.
    /// Used internally by copy/move operations to write to the destination.
    pub fn as_destination(&self) -> Client {
        let inner = Arc::new(ClientInner {
            src: Side {
                logical: self.inner.dst.logical.clone(),
                mount_provider: self.inner.dst.mount_provider.clone(),
            },
            dst: Side {
                logical: self.inner.dst.logical.clone(),
                mount_provider: self.inner.dst.mount_provider.clone(),
            },
            workers: self.inner.workers,
            absolute_path: self.inner.absolute_path,
            ignore_access_errors: self.inner.ignore_access_errors,
        });
        Client { inner }
    }

    /// Rewrite paths in `data` from `src` to `dst`, honouring `absolute_path`.
    pub fn swap_paths(
        &self,
        data: &mut BTreeMap<String, Map<String, Value>>,
        src: &str,
        dst: &str,
    ) {
        if self.absolute_path() {
            trim_prefix_map(data, src);
        }
        ensure_prefix_map(data, dst);
    }

    /// Prepare a path for input into a read/write/list/delete given the user's
    /// `absolute_path` preference.
    pub fn input_path(&self, path: &str, root: &str) -> String {
        if self.absolute_path() {
            path.to_string()
        } else {
            crate::api::helpers::add_prefix(path, root)
        }
    }

    /// Prepare a path for output to the user.
    pub fn output_path(&self, path: &str, root: &str) -> String {
        if self.absolute_path() {
            crate::api::helpers::ensure_prefix(path, root)
        } else {
            path_join(&[path.strip_prefix(root).unwrap_or(path)])
        }
    }

    /// Prepare a list of paths for output to the user.
    pub fn output_paths(&self, paths: &mut [String], root: &str) {
        if self.absolute_path() {
            add_prefix_list(paths, root);
        }
    }
}

/// Builder for [`Client`].
#[derive(Default)]
pub struct ClientBuilder {
    src_logical: Option<Arc<dyn Logical>>,
    dst_logical: Option<Arc<dyn Logical>>,
    src_mount_provider: Option<Arc<dyn MountProvider>>,
    dst_mount_provider: Option<Arc<dyn MountProvider>>,
    workers: Option<usize>,
    absolute_path: bool,
    ignore_access_errors: bool,
}

impl ClientBuilder {
    /// Set both the source-side logical (`WithVaultClient`).
    pub fn with_vault_client(mut self, client: VaultHttpClient) -> Self {
        self.src_logical = Some(Arc::new(client));
        self
    }

    /// Set the source-side logical via any [`Logical`] impl.
    pub fn with_logical(mut self, logical: Arc<dyn Logical>) -> Self {
        self.src_logical = Some(logical);
        self
    }

    /// Alias of [`with_vault_client`].
    pub fn with_vault_src_client(self, client: VaultHttpClient) -> Self {
        self.with_vault_client(client)
    }

    /// Set the destination-side Vault client (`WithVaultDstClient`).
    pub fn with_vault_dst_client(mut self, client: VaultHttpClient) -> Self {
        self.dst_logical = Some(Arc::new(client));
        self
    }

    /// Set the destination-side logical via any [`Logical`] impl.
    pub fn with_dst_logical(mut self, logical: Arc<dyn Logical>) -> Self {
        self.dst_logical = Some(logical);
        self
    }

    /// Set the worker concurrency limit.
    pub fn with_workers(mut self, n: usize) -> Self {
        self.workers = Some(n);
        self
    }

    /// Toggle absolute-path output mode.
    pub fn with_absolute_path(mut self, b: bool) -> Self {
        self.absolute_path = b;
        self
    }

    /// Toggle silent ignore of read/list access errors.
    pub fn with_ignore_access_errors(mut self, b: bool) -> Self {
        self.ignore_access_errors = b;
        self
    }

    /// Set the source-side mount provider (alias of `with_src_mount_provider`).
    pub fn with_mount_provider(mut self, p: Arc<dyn MountProvider>) -> Self {
        self.src_mount_provider = Some(p);
        self
    }

    /// Set the source-side mount provider.
    pub fn with_src_mount_provider(self, p: Arc<dyn MountProvider>) -> Self {
        self.with_mount_provider(p)
    }

    /// Set the destination-side mount provider.
    pub fn with_dst_mount_provider(mut self, p: Arc<dyn MountProvider>) -> Self {
        self.dst_mount_provider = Some(p);
        self
    }

    /// Construct the client.
    pub fn build(self) -> Result<Client, Error> {
        let workers = self.workers.unwrap_or(DEFAULT_WORKERS);
        if workers < 1 {
            let msg = format!("workers must 1 or greater: {workers}");
            return Err(Error::wrap(
                "",
                ErrorKind::ApplyOptions,
                Some(Box::new(Error::wrap(&msg, ErrorKind::NumWorkers, None))),
            ));
        }

        let src_logical = self.src_logical.unwrap_or_else(|| Arc::new(NullLogical));
        let dst_logical = self.dst_logical.unwrap_or_else(|| src_logical.clone());

        let src_mount_provider = self.src_mount_provider.unwrap_or_else(|| {
            Arc::new(DefaultMountProvider {
                logical: src_logical.clone(),
            })
        });
        let dst_mount_provider = self.dst_mount_provider.unwrap_or_else(|| {
            Arc::new(DefaultMountProvider {
                logical: dst_logical.clone(),
            })
        });

        Ok(Client {
            inner: Arc::new(ClientInner {
                src: Side {
                    logical: src_logical,
                    mount_provider: src_mount_provider,
                },
                dst: Side {
                    logical: dst_logical,
                    mount_provider: dst_mount_provider,
                },
                workers,
                absolute_path: self.absolute_path,
                ignore_access_errors: self.ignore_access_errors,
            }),
        })
    }
}

/// A no-op [`Logical`] used when the builder produced a client without one
/// (so trait objects always have a concrete impl). All methods return an
/// error so that any unexpected use surfaces clearly.
#[derive(Debug)]
struct NullLogical;

#[async_trait::async_trait]
impl Logical for NullLogical {
    async fn read(&self, _: &str) -> Result<Option<crate::api::logical::Secret>, crate::api::error::BoxError> {
        Err("no vault client configured".into())
    }
    async fn read_with_data(
        &self,
        _: &str,
        _: &[(&str, &str)],
    ) -> Result<Option<crate::api::logical::Secret>, crate::api::error::BoxError> {
        Err("no vault client configured".into())
    }
    async fn list(&self, _: &str) -> Result<Option<crate::api::logical::Secret>, crate::api::error::BoxError> {
        Err("no vault client configured".into())
    }
    async fn write(&self, _: &str, _: serde_json::Value) -> Result<Option<crate::api::logical::Secret>, crate::api::error::BoxError> {
        Err("no vault client configured".into())
    }
    async fn delete(&self, _: &str) -> Result<Option<crate::api::logical::Secret>, crate::api::error::BoxError> {
        Err("no vault client configured".into())
    }
}

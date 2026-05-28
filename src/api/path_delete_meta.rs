//! `PathDeleteMeta` — wipe all metadata and versions of a KV v2 secret.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::mounts::VaultOp;

impl Client {
    /// Delete all secret metadata and versions at `p`. KV v2 only.
    pub async fn path_delete_meta(&self, p: &str) -> Result<(), Error> {
        self.path_delete_with_op(p, VaultOp::DeleteMeta)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathDeleteMeta, Some(Box::new(e))))
    }
}

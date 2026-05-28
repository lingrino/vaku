//! `PathDelete`.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::mounts::{rewrite_path, VaultOp};

impl Client {
    /// Delete the current version at `p`. On KV v2 this is a soft delete.
    pub async fn path_delete(&self, p: &str) -> Result<(), Error> {
        self.path_delete_with_op(p, VaultOp::Delete)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathDelete, Some(Box::new(e))))
    }

    pub(crate) async fn path_delete_with_op(&self, p: &str, op: VaultOp) -> Result<(), Error> {
        let (vault_path, _) =
            rewrite_path(self.src().mount_provider.as_ref(), p, op).await?;

        self.src().logical.delete(&vault_path).await.map_err(|e| {
            Error::wrap(&e.to_string(), ErrorKind::VaultDelete, None)
        })?;
        Ok(())
    }
}

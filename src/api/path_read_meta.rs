//! `PathReadMeta` — fetch the metadata blob for a secret (KV v2 only).

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::mounts::{rewrite_path, VaultOp};
use crate::api::secret::{extract_secret_meta, SecretMeta};

impl Client {
    /// Read all metadata for a KV v2 secret. Returns `None` for missing
    /// secrets; an error for KV v1 mounts.
    pub async fn path_read_meta(&self, p: &str) -> Result<Option<SecretMeta>, Error> {
        let (vault_path, _) =
            rewrite_path(self.src().mount_provider.as_ref(), p, VaultOp::ReadMeta)
                .await
                .map_err(|e| Error::wrap(p, ErrorKind::PathReadMeta, Some(Box::new(e))))?;

        let secret = self.src().logical.read(&vault_path).await.map_err(|e| {
            Error::wrap(
                p,
                ErrorKind::PathReadMeta,
                Some(Box::new(Error::wrap(
                    &e.to_string(),
                    ErrorKind::VaultRead,
                    None,
                ))),
            )
        })?;

        let Some(secret) = secret else {
            return Ok(None);
        };
        let Some(data) = secret.data else {
            return Ok(None);
        };
        Ok(Some(extract_secret_meta(Some(&data))))
    }
}

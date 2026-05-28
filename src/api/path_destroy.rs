//! `PathDestroy` — permanently destroy specific versions of a KV v2 secret.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::mounts::{rewrite_path, VaultOp};
use crate::api::secret::KV2_VERSIONS;
use serde_json::{json, Map, Value};

impl Client {
    /// Destroy the given `versions` of the secret at `p`. KV v2 only.
    pub async fn path_destroy(&self, p: &str, versions: &[i64]) -> Result<(), Error> {
        if versions.is_empty() {
            return Err(Error::wrap(
                "no versions provided",
                ErrorKind::PathDestroy,
                None,
            ));
        }

        let (vault_path, _) = rewrite_path(self.src().mount_provider.as_ref(), p, VaultOp::Destroy)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathDestroy, Some(Box::new(e))))?;

        let mut body = Map::new();
        body.insert(KV2_VERSIONS.to_string(), json!(versions.to_vec()));

        self.src()
            .logical
            .write(&vault_path, Value::Object(body))
            .await
            .map_err(|e| {
                Error::wrap(
                    p,
                    ErrorKind::PathDestroy,
                    Some(Box::new(Error::wrap(
                        &e.to_string(),
                        ErrorKind::VaultWrite,
                        None,
                    ))),
                )
            })?;
        Ok(())
    }
}

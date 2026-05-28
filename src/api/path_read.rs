//! `PathRead` and `PathReadVersion`.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::mounts::{rewrite_path, MountVersion, VaultOp};
use crate::api::secret::{extract_v2_read, KV2_VERSION};
use serde_json::{Map, Value};

impl Client {
    /// Read the current data at a path. Returns `None` when the path doesn't
    /// exist or has been deleted/destroyed.
    pub async fn path_read(&self, p: &str) -> Result<Option<Map<String, Value>>, Error> {
        let (vault_path, mv) = rewrite_path(self.src().mount_provider.as_ref(), p, VaultOp::Read)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathRead, Some(Box::new(e))))?;

        let res = self.src().logical.read(&vault_path).await;
        let secret = match res {
            Ok(s) => s,
            Err(e) if self.ignore_access_errors() => {
                let _ = e;
                return Ok(None);
            }
            Err(e) => {
                return Err(Error::wrap(
                    p,
                    ErrorKind::PathRead,
                    Some(Box::new(Error::wrap(&e.to_string(), ErrorKind::VaultRead, None))),
                ))
            }
        };

        let Some(secret) = secret else { return Ok(None) };
        let Some(data) = secret.data else { return Ok(None) };

        if matches!(mv, MountVersion::Mv2) {
            Ok(extract_v2_read(Some(&data)))
        } else {
            Ok(Some(data))
        }
    }

    /// Read a specific version of a secret. KV v2 only.
    pub async fn path_read_version(
        &self,
        p: &str,
        version: i64,
    ) -> Result<Option<Map<String, Value>>, Error> {
        let (vault_path, mv) = rewrite_path(self.src().mount_provider.as_ref(), p, VaultOp::Read)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathReadVersion, Some(Box::new(e))))?;

        if !matches!(mv, MountVersion::Mv2) {
            return Err(Error::wrap(
                p,
                ErrorKind::PathReadVersion,
                Some(Box::new(Error::wrap("", ErrorKind::MountVersion, None))),
            ));
        }

        let version_s = version.to_string();
        let secret = self
            .src()
            .logical
            .read_with_data(&vault_path, &[(KV2_VERSION, &version_s)])
            .await
            .map_err(|e| {
                Error::wrap(
                    p,
                    ErrorKind::PathReadVersion,
                    Some(Box::new(Error::wrap(&e.to_string(), ErrorKind::VaultRead, None))),
                )
            })?;

        let Some(secret) = secret else { return Ok(None) };
        let Some(data) = secret.data else { return Ok(None) };
        Ok(extract_v2_read(Some(&data)))
    }
}

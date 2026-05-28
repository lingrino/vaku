//! `PathList`: list a single path's immediate children.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::logical::Secret;
use crate::api::mounts::{rewrite_path, VaultOp};
use serde_json::Value;

impl Client {
    /// List the immediate child paths under `p`.
    pub async fn path_list(&self, p: &str) -> Result<Vec<String>, Error> {
        let (vault_path, _) = rewrite_path(self.src().mount_provider.as_ref(), p, VaultOp::List)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathList, Some(Box::new(e))))?;

        let secret = match self.src().logical.list(&vault_path).await {
            Ok(s) => s,
            Err(_) if self.ignore_access_errors() => return Ok(Vec::new()),
            Err(e) => {
                return Err(Error::wrap(
                    p,
                    ErrorKind::PathList,
                    Some(Box::new(Error::wrap(
                        &e.to_string(),
                        ErrorKind::VaultList,
                        None,
                    ))),
                ))
            }
        };

        let mut list = decode_secret_keys(secret.as_ref())
            .map_err(|e| Error::wrap(p, ErrorKind::PathList, Some(Box::new(e))))?;

        self.output_paths(&mut list, p);
        Ok(list)
    }
}

/// Decode the `keys` array from a Vault LIST response.
pub(crate) fn decode_secret_keys(secret: Option<&Secret>) -> Result<Vec<String>, Error> {
    let Some(secret) = secret else {
        return Ok(Vec::new());
    };
    let Some(data) = &secret.data else {
        return Ok(Vec::new());
    };

    let raw = match data.get("keys") {
        Some(v) => v,
        None => return Err(Error::wrap("", ErrorKind::DecodeSecret, None)),
    };
    if raw.is_null() {
        return Err(Error::wrap("", ErrorKind::DecodeSecret, None));
    }
    let arr = match raw.as_array() {
        Some(a) => a,
        None => return Err(Error::wrap("", ErrorKind::DecodeSecret, None)),
    };

    let mut out = Vec::with_capacity(arr.len());
    for item in arr {
        match item {
            Value::String(s) => out.push(s.clone()),
            _ => return Err(Error::wrap("", ErrorKind::DecodeSecret, None)),
        }
    }
    Ok(out)
}

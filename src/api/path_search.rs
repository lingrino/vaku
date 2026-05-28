//! `PathSearch` — text-search a single secret.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use serde_json::{Map, Value};

impl Client {
    /// Returns true if any key in the secret at `p` contains `search`, or if
    /// the JSON representation of any value does. Mirrors Go's `PathSearch`.
    pub async fn path_search(&self, p: &str, search: &str) -> Result<bool, Error> {
        let read = self
            .path_read(p)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathSearch, Some(Box::new(e))))?;
        let Some(secret) = read else { return Ok(false) };
        search_secret(&secret, search)
            .map_err(|e| Error::wrap(p, ErrorKind::PathSearch, Some(Box::new(e))))
    }
}

/// Shared by `folder_search`. Returns `Ok(false)` for no match; an error only
/// if serialization fails.
pub(crate) fn search_secret(secret: &Map<String, Value>, search: &str) -> Result<bool, Error> {
    for (k, v) in secret {
        if k.contains(search) {
            return Ok(true);
        }
        let s = serde_json::to_string(v)
            .map_err(|_| Error::wrap("", ErrorKind::JsonMarshal, None))?;
        if s.contains(search) {
            return Ok(true);
        }
    }
    Ok(false)
}

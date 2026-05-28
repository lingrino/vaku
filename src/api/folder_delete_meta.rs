//! `FolderDeleteMeta` — recursive `path_delete_meta` (KV v2 only).

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use std::sync::Arc;

impl Client {
    /// Recursively wipe every secret + its metadata under `p`. KV v2 only.
    pub async fn folder_delete_meta(&self, p: &str) -> Result<(), Error> {
        let deleter = Arc::new(|c: Client, path: String| {
            Box::pin(async move { c.path_delete_meta(&path).await })
                as std::pin::Pin<Box<dyn std::future::Future<Output = Result<(), Error>> + Send>>
        });
        self.folder_delete_with(p, deleter, ErrorKind::FolderDeleteMeta).await
    }
}

//! `FolderSearch` — recursive text search.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::path_search::search_secret;

impl Client {
    /// Recursively read every secret under `path` and return the paths whose
    /// keys or values contain `search`.
    pub async fn folder_search(&self, path: &str, search: &str) -> Result<Vec<String>, Error> {
        let read = self
            .folder_read(path)
            .await
            .map_err(|e| Error::wrap(path, ErrorKind::FolderSearch, Some(Box::new(e))))?;

        let Some(read) = read else { return Ok(Vec::new()) };

        let mut matches = Vec::new();
        for (p, secret) in read {
            let found = search_secret(&secret, search)
                .map_err(|e| Error::wrap(path, ErrorKind::FolderSearch, Some(Box::new(e))))?;
            if found {
                matches.push(p);
            }
        }
        Ok(matches)
    }
}

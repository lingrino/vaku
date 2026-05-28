//! `FolderMoveAllVersions` — folder_copy_all_versions + folder_delete_meta.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};

impl Client {
    /// Recursive all-versions move (copy then wipe metadata). KV v2 only.
    pub async fn folder_move_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        self.folder_copy_all_versions(src, dst)
            .await
            .map_err(|e| Error::wrap("", ErrorKind::FolderMoveAllVersions, Some(Box::new(e))))?;
        self.folder_delete_meta(src)
            .await
            .map_err(|e| Error::wrap(src, ErrorKind::FolderMoveAllVersions, Some(Box::new(e))))
    }
}

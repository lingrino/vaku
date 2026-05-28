//! `FolderMove` — folder_copy + folder_delete on source.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};

impl Client {
    /// Recursive move (copy + delete) from `src` to `dst`.
    pub async fn folder_move(&self, src: &str, dst: &str) -> Result<(), Error> {
        self.folder_copy(src, dst)
            .await
            .map_err(|e| Error::wrap("", ErrorKind::FolderMove, Some(Box::new(e))))?;
        self.folder_delete(src).await.map_err(|e| {
            Error::wrap(
                &format!("delete {src}"),
                ErrorKind::FolderMove,
                Some(Box::new(e)),
            )
        })
    }
}

//! `FolderCopy` — recursive folder copy via read + swap_paths + dst write.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};

impl Client {
    /// Copy every secret under `src` to `dst`. Cross-cluster safe.
    pub async fn folder_copy(&self, src: &str, dst: &str) -> Result<(), Error> {
        let mut read = self
            .folder_read(src)
            .await
            .map_err(|e| {
                Error::wrap(
                    &format!("read from {src}"),
                    ErrorKind::FolderCopy,
                    Some(Box::new(e)),
                )
            })?
            .unwrap_or_default();

        if read.is_empty() {
            return Ok(());
        }

        self.swap_paths(&mut read, src, dst);

        let dst_client = self.as_destination();
        dst_client.folder_write(read).await.map_err(|e| {
            Error::wrap(
                &format!("write to {dst}"),
                ErrorKind::FolderCopy,
                Some(Box::new(e)),
            )
        })
    }
}

//! `PathMoveAllVersions` — copy all versions then delete source metadata.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};

impl Client {
    /// Move every version of `src` to `dst`. KV v2 only.
    pub async fn path_move_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        self.path_copy_all_versions(src, dst)
            .await
            .map_err(|e| Error::wrap("", ErrorKind::PathMoveAllVersions, Some(Box::new(e))))?;
        self.path_delete_meta(src)
            .await
            .map_err(|e| Error::wrap(src, ErrorKind::PathMoveAllVersions, Some(Box::new(e))))
    }
}

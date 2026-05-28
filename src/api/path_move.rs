//! `PathMove` — copy then delete source.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};

impl Client {
    /// Move (copy + delete) the current version from `src` to `dst`.
    pub async fn path_move(&self, src: &str, dst: &str) -> Result<(), Error> {
        self.path_copy(src, dst)
            .await
            .map_err(|e| Error::wrap("", ErrorKind::PathMove, Some(Box::new(e))))?;
        self.path_delete(src)
            .await
            .map_err(|e| Error::wrap(dst, ErrorKind::PathMove, Some(Box::new(e))))
    }
}

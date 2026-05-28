//! `PathCopy` — read from src, write to dst.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};

impl Client {
    /// Copy the current secret from `src` to `dst` (possibly across vault
    /// clusters).
    pub async fn path_copy(&self, src: &str, dst: &str) -> Result<(), Error> {
        let secret = self
            .path_read(src)
            .await
            .map_err(|e| Error::wrap(src, ErrorKind::PathCopy, Some(Box::new(e))))?;

        // Mirror Go: passing `nil` to PathWrite is an error (ErrNilData). The
        // destination client is the dst-side view of self.
        let dst_client = self.as_destination();
        dst_client
            .path_write(dst, secret)
            .await
            .map_err(|e| Error::wrap(dst, ErrorKind::PathCopy, Some(Box::new(e))))
    }
}

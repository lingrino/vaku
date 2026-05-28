//! `PathWrite`.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::mounts::{rewrite_path, MountVersion, VaultOp};
use crate::api::secret::KV2_DATA;
use serde_json::{Map, Value};

impl Client {
    /// Write `data` to `p`. Mirrors Go's `PathWrite` — a `None` data argument
    /// yields [`ErrorKind::NilData`].
    pub async fn path_write(
        &self,
        p: &str,
        data: Option<Map<String, Value>>,
    ) -> Result<(), Error> {
        let Some(data) = data else {
            return Err(Error::wrap(
                p,
                ErrorKind::PathWrite,
                Some(Box::new(Error::wrap("", ErrorKind::NilData, None))),
            ));
        };

        let (vault_path, mv) = rewrite_path(self.src().mount_provider.as_ref(), p, VaultOp::Write)
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathWrite, Some(Box::new(e))))?;

        let body: Value = if matches!(mv, MountVersion::Mv2) {
            let mut wrapped = Map::new();
            wrapped.insert(KV2_DATA.to_string(), Value::Object(data));
            Value::Object(wrapped)
        } else {
            Value::Object(data)
        };

        self.src().logical.write(&vault_path, body).await.map_err(|e| {
            Error::wrap(
                p,
                ErrorKind::PathWrite,
                Some(Box::new(Error::wrap(&e.to_string(), ErrorKind::VaultWrite, None))),
            )
        })?;
        Ok(())
    }
}

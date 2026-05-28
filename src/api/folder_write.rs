//! `FolderWrite` — write many secrets concurrently.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use serde_json::{Map, Value};
use std::collections::BTreeMap;
use std::sync::Arc;
use tokio::task::JoinSet;

impl Client {
    /// Write every entry of `data` to its corresponding path.
    pub async fn folder_write(
        &self,
        data: BTreeMap<String, Map<String, Value>>,
    ) -> Result<(), Error> {
        let data = Arc::new(data);

        let (tx, rx) = async_channel::unbounded::<String>();
        for k in data.keys() {
            let _ = tx.try_send(k.clone());
        }
        drop(tx);

        let mut set = JoinSet::new();
        for _ in 0..self.workers() {
            let client = self.clone();
            let rx = rx.clone();
            let data = data.clone();
            set.spawn(async move {
                while let Ok(path) = rx.recv().await {
                    let inner = data.get(&path).cloned();
                    client.path_write(&path, inner).await?;
                }
                Ok::<(), Error>(())
            });
        }

        let mut first_err: Option<Error> = None;
        while let Some(r) = set.join_next().await {
            match r {
                Ok(Ok(())) => {}
                Ok(Err(e)) => {
                    if first_err.is_none() {
                        first_err = Some(e);
                    }
                }
                Err(je) => {
                    if first_err.is_none() {
                        first_err = Some(Error::from_msg(je.to_string()));
                    }
                }
            }
        }
        if let Some(e) = first_err {
            return Err(Error::wrap("", ErrorKind::FolderWrite, Some(Box::new(e))));
        }
        Ok(())
    }
}

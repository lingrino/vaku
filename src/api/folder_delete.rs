//! `FolderDelete` — recursively delete every secret under a folder.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use std::future::Future;
use std::pin::Pin;
use std::sync::Arc;
use tokio::task::JoinSet;

/// Boxed async closure for one of (delete | delete_meta).
type Deleter = Arc<
    dyn (Fn(Client, String) -> Pin<Box<dyn Future<Output = Result<(), Error>> + Send>>)
        + Send
        + Sync,
>;

impl Client {
    /// Recursively delete every secret under `p`. Soft delete on KV v2.
    pub async fn folder_delete(&self, p: &str) -> Result<(), Error> {
        let deleter: Deleter =
            Arc::new(|c: Client, path: String| Box::pin(async move { c.path_delete(&path).await }));
        self.folder_delete_with(p, deleter, ErrorKind::FolderDelete)
            .await
    }

    pub(crate) async fn folder_delete_with(
        &self,
        p: &str,
        deleter: Deleter,
        outer_kind: ErrorKind,
    ) -> Result<(), Error> {
        let stream = self.folder_list_chan(p);
        let mut done = stream.done;
        let path_rx = stream.results;

        let mut set: JoinSet<Result<(), Error>> = JoinSet::new();
        for _ in 0..self.workers() {
            let client = self.clone();
            let rx = path_rx.clone();
            let root = p.to_string();
            let deleter = deleter.clone();
            set.spawn(async move {
                while let Ok(path) = rx.recv().await {
                    let in_path = client.input_path(&path, &root);
                    (deleter)(client.clone(), in_path).await?;
                }
                Ok::<(), Error>(())
            });
        }

        let mut first_err: Option<Error> = None;
        let mut workers_done = false;
        let mut list_done = false;
        while !(workers_done && list_done) {
            tokio::select! {
                biased;
                rd = set.join_next(), if !workers_done => match rd {
                    Some(Ok(Ok(()))) => {}
                    Some(Ok(Err(e))) => if first_err.is_none() { first_err = Some(e) },
                    Some(Err(je)) => if first_err.is_none() {
                        first_err = Some(Error::from_msg(je.to_string()));
                    },
                    None => workers_done = true,
                },
                d = &mut done, if !list_done => {
                    list_done = true;
                    if let Ok(Err(e)) = d {
                        if first_err.is_none() { first_err = Some(e); }
                    }
                }
            }
        }

        if let Some(e) = first_err {
            return Err(Error::wrap(p, outer_kind, Some(Box::new(e))));
        }
        Ok(())
    }
}

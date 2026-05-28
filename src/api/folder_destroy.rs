//! `FolderDestroy` — destroy listed versions of every secret in a folder.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use tokio::task::JoinSet;

impl Client {
    /// Destroy `versions` of every secret under `p`. KV v2 only.
    pub async fn folder_destroy(&self, p: &str, versions: &[i64]) -> Result<(), Error> {
        let versions: Vec<i64> = versions.to_vec();

        let stream = self.folder_list_chan(p);
        let mut done = stream.done;
        let path_rx = stream.results;

        let mut set: JoinSet<Result<(), Error>> = JoinSet::new();
        for _ in 0..self.workers() {
            let client = self.clone();
            let rx = path_rx.clone();
            let root = p.to_string();
            let versions = versions.clone();
            set.spawn(async move {
                while let Ok(path) = rx.recv().await {
                    let in_path = client.input_path(&path, &root);
                    client.path_destroy(&in_path, &versions).await?;
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
            return Err(Error::wrap(p, ErrorKind::FolderDestroy, Some(Box::new(e))));
        }
        Ok(())
    }
}

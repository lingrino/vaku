//! `FolderRead` — recursive secret read.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::helpers::merge_maps;
use serde_json::{Map, Value};
use std::collections::BTreeMap;
use tokio::task::JoinSet;

pub type FolderReadResult = BTreeMap<String, Map<String, Value>>;

impl Client {
    /// Recursively read every secret under `p`.
    pub async fn folder_read(&self, p: &str) -> Result<Option<FolderReadResult>, Error> {
        let stream = self.folder_list_chan(p);
        let mut done = stream.done;
        let path_rx = stream.results;

        let (res_tx, mut res_rx) =
            tokio::sync::mpsc::unbounded_channel::<(String, Map<String, Value>)>();

        // Spawn `workers` readers; each drains the path stream into res_tx.
        let mut workers: JoinSet<Result<(), Error>> = JoinSet::new();
        for _ in 0..self.workers() {
            let client = self.clone();
            let rx = path_rx.clone();
            let root = p.to_string();
            let res_tx = res_tx.clone();
            workers.spawn(async move {
                while let Ok(path) = rx.recv().await {
                    let in_path = client.input_path(&path, &root);
                    match client.path_read(&in_path).await {
                        Ok(None) => {}
                        Ok(Some(data)) => {
                            let key = client.output_path(&in_path, &root);
                            // Receiver was dropped: caller errored out; bail.
                            if res_tx.send((key, data)).is_err() {
                                return Ok(());
                            }
                        }
                        Err(e) => return Err(e),
                    }
                }
                Ok(())
            });
        }
        drop(res_tx);

        let mut out: FolderReadResult = BTreeMap::new();
        let mut list_err: Option<Error> = None;
        let mut worker_err: Option<Error> = None;
        let mut list_done = false;
        let mut workers_done = false;

        while !(list_done && workers_done && res_rx.is_empty()) {
            tokio::select! {
                biased;
                Some((k, v)) = res_rx.recv(), if !res_rx.is_closed() || !res_rx.is_empty() => {
                    out.insert(k, v);
                }
                rd = workers.join_next(), if !workers_done => {
                    match rd {
                        Some(Ok(Ok(()))) => {}
                        Some(Ok(Err(e))) => {
                            if worker_err.is_none() { worker_err = Some(e); }
                        }
                        Some(Err(je)) => {
                            if worker_err.is_none() {
                                worker_err = Some(Error::from_msg(je.to_string()));
                            }
                        }
                        None => { workers_done = true; }
                    }
                }
                d = &mut done, if !list_done => {
                    list_done = true;
                    if let Ok(Err(e)) = d { list_err = Some(e); }
                }
            }
        }

        if let Some(e) = list_err {
            return Err(Error::wrap(
                p,
                ErrorKind::FolderRead,
                Some(Box::new(Error::wrap(
                    p,
                    ErrorKind::FolderReadChan,
                    Some(Box::new(e)),
                ))),
            ));
        }
        if let Some(e) = worker_err {
            return Err(Error::wrap(
                p,
                ErrorKind::FolderRead,
                Some(Box::new(Error::wrap(
                    p,
                    ErrorKind::FolderReadChan,
                    Some(Box::new(e)),
                ))),
            ));
        }

        Ok(if out.is_empty() { None } else { Some(out) })
    }
}

/// Public helper exposed for callers that want to merge a folder-read result.
pub fn merge(into: &mut FolderReadResult, from: FolderReadResult) {
    merge_maps(into, from)
}

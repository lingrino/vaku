//! `FolderList` and `FolderListChan` — recursive listing.
//!
//! This is the recursion engine that every other folder operation builds on.
//! It mirrors Go's `FolderListChan`: workers pull from `path_rx`, and when
//! they see a folder they re-enqueue its children back into `path_tx`. We
//! track outstanding paths via an atomic counter so the worker pool can
//! cleanly drain when there's no more work.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::helpers::{ensure_folder, is_folder};
use async_channel::{Receiver, Sender};
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use tokio::sync::oneshot;
use tokio_util::sync::CancellationToken;

/// Streaming results from [`Client::folder_list_chan`].
pub struct FolderListStream {
    pub results: Receiver<String>,
    pub done: oneshot::Receiver<Result<(), Error>>,
}

impl Client {
    /// Recursively list every leaf path under `p`. The order of results is
    /// unspecified (parallel workers).
    pub async fn folder_list(&self, p: &str) -> Result<Vec<String>, Error> {
        let stream = self.folder_list_chan(p);
        let mut out = Vec::new();
        // Race the results stream against the completion signal so we exit
        // promptly on error.
        let mut done = stream.done;
        let results = stream.results;
        loop {
            tokio::select! {
                biased;
                msg = results.recv() => match msg {
                    Ok(s) => out.push(s),
                    Err(_) => {
                        // Channel closed -> work finished -> await done signal.
                        match (&mut done).await {
                            Ok(Ok(())) => return Ok(out),
                            Ok(Err(e)) => return Err(Error::wrap(p, ErrorKind::FolderList, Some(Box::new(e)))),
                            Err(_) => return Ok(out),
                        }
                    }
                },
                d = &mut done => {
                    if let Ok(Err(e)) = d {
                        return Err(Error::wrap(p, ErrorKind::FolderList, Some(Box::new(e))));
                    }
                    // Drain any remaining results before returning.
                    while let Ok(s) = results.try_recv() {
                        out.push(s);
                    }
                    return Ok(out);
                }
            }
        }
    }

    /// Returns streams that incrementally produce listed paths and a final
    /// completion signal. Mirrors Go's `FolderListChan`.
    pub fn folder_list_chan(&self, p: &str) -> FolderListStream {
        let client = self.clone();
        let workers = client.workers();
        let root = ensure_folder(p);

        let (path_tx, path_rx) = async_channel::unbounded::<String>();
        let (res_tx, res_rx) = async_channel::unbounded::<String>();
        let (done_tx, done_rx) = oneshot::channel::<Result<(), Error>>();

        let pending = Arc::new(AtomicUsize::new(1));
        let cancel = CancellationToken::new();
        let _ = path_tx.try_send(root.clone());

        let mut handles = Vec::with_capacity(workers);
        for _ in 0..workers {
            let client = client.clone();
            let root = root.clone();
            let path_tx = path_tx.clone();
            let path_rx = path_rx.clone();
            let res_tx = res_tx.clone();
            let pending = pending.clone();
            let cancel = cancel.clone();
            handles.push(tokio::spawn(async move {
                folder_list_worker(client, root, path_tx, path_rx, res_tx, pending, cancel).await
            }));
        }

        // Supervisor: collect worker results; first error -> cancel everyone
        // and signal done.
        let supervisor_path_tx = path_tx.clone();
        let supervisor_res_tx = res_tx.clone();
        tokio::spawn(async move {
            let mut first_err: Option<Error> = None;
            for h in handles {
                match h.await {
                    Ok(Ok(())) => {}
                    Ok(Err(e)) => {
                        if first_err.is_none() {
                            first_err = Some(e);
                            cancel.cancel();
                        }
                    }
                    Err(join_err) => {
                        if first_err.is_none() {
                            first_err = Some(Error::wrap(
                                "worker panic",
                                ErrorKind::Custom(join_err.to_string()),
                                None,
                            ));
                            cancel.cancel();
                        }
                    }
                }
            }
            // Close channels so consumers stop blocking.
            drop(supervisor_path_tx);
            drop(supervisor_res_tx);
            let _ = done_tx.send(match first_err {
                Some(e) => Err(e),
                None => Ok(()),
            });
        });

        FolderListStream {
            results: res_rx,
            done: done_rx,
        }
    }
}

async fn folder_list_worker(
    client: Client,
    root: String,
    path_tx: Sender<String>,
    path_rx: Receiver<String>,
    res_tx: Sender<String>,
    pending: Arc<AtomicUsize>,
    cancel: CancellationToken,
) -> Result<(), Error> {
    loop {
        tokio::select! {
            biased;
            _ = cancel.cancelled() => {
                return Err(Error::wrap("", ErrorKind::Context, Some(Box::new(Error::from_msg("cancelled")))));
            }
            msg = path_rx.recv() => {
                match msg {
                    Err(_) => return Ok(()),
                    Ok(path) => {
                        let r = process_path(&client, &root, &path, &path_tx, &res_tx, &pending).await;
                        // Always decrement pending for the path we just consumed.
                        let prev = pending.fetch_sub(1, Ordering::SeqCst);
                        if prev == 1 {
                            // No more outstanding paths — close the queue so
                            // every worker can exit.
                            path_tx.close();
                            res_tx.close();
                        }
                        r?;
                    }
                }
            }
        }
    }
}

async fn process_path(
    client: &Client,
    root: &str,
    path: &str,
    path_tx: &Sender<String>,
    res_tx: &Sender<String>,
    pending: &Arc<AtomicUsize>,
) -> Result<(), Error> {
    if is_folder(path) {
        let list = client
            .path_list(path)
            .await
            .map_err(|e| Error::wrap(root, ErrorKind::FolderListChan, Some(Box::new(e))))?;
        // Increment pending BEFORE sending so the count can't transiently
        // reach zero between consumption and enqueue.
        if !list.is_empty() {
            pending.fetch_add(list.len(), Ordering::SeqCst);
            for child in list {
                let next = client.input_path(&child, path);
                if path_tx.send(next).await.is_err() {
                    // Channel closed (cancellation): roll back the count.
                    pending.fetch_sub(1, Ordering::SeqCst);
                }
            }
        }
    } else {
        let out = client.output_path(path, root);
        let _ = res_tx.send(out).await;
    }
    Ok(())
}

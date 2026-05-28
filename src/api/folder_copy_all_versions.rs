//! `FolderCopyAllVersions` — per-leaf path_copy_all_versions across a folder.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::helpers::path_join;
use crate::api::path_copy_all_versions::validate_kv2;
use tokio::task::JoinSet;

impl Client {
    /// Recursively copy every version of every secret under `src` to `dst`.
    /// KV v2 only for both sides.
    pub async fn folder_copy_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        validate_kv2(self.src().mount_provider.as_ref(), src)
            .await
            .map_err(|e| Error::wrap(src, ErrorKind::FolderCopyAllVersions, Some(Box::new(e))))?;
        validate_kv2(self.dst().mount_provider.as_ref(), dst)
            .await
            .map_err(|e| Error::wrap(dst, ErrorKind::FolderCopyAllVersions, Some(Box::new(e))))?;

        let stream = self.folder_list_chan(src);
        let mut done = stream.done;
        let path_rx = stream.results;

        let mut set: JoinSet<Result<(), Error>> = JoinSet::new();
        let abs = self.absolute_path();
        for _ in 0..self.workers() {
            let client = self.clone();
            let rx = path_rx.clone();
            let src_root = src.to_string();
            let dst_root = dst.to_string();
            set.spawn(async move {
                while let Ok(path) = rx.recv().await {
                    let src_path = client.input_path(&path, &src_root);

                    // Mirror Go's branching for `absolute_path`:
                    //   abs:   strip src prefix, then prepend dst
                    //   rel:   prepend dst (path is already relative)
                    let dst_path = if abs {
                        let rel = path.strip_prefix(&src_root as &str).unwrap_or(&path);
                        path_join(&[&dst_root, rel])
                    } else {
                        client.input_path(&path, &dst_root)
                    };

                    client.path_copy_all_versions(&src_path, &dst_path).await?;
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
            return Err(Error::wrap(
                src,
                ErrorKind::FolderCopyAllVersions,
                Some(Box::new(e)),
            ));
        }
        Ok(())
    }
}

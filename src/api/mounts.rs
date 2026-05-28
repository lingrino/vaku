//! Mount-info lookup and KV-v2 path rewriting.

use crate::api::error::{Error, ErrorKind};
use crate::api::helpers::{ensure_folder, insert_into_path};
use crate::api::secret::{KV2_DATA, KV2_DESTROY, KV2_METADATA};

/// KV mount versions recognized by Vaku.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum MountVersion {
    Mv0,
    Mv1,
    Mv2,
    Other(i64),
}

impl MountVersion {
    pub fn from_str(s: &str) -> Self {
        match s.parse::<i64>() {
            Ok(0) => MountVersion::Mv0,
            Ok(1) => MountVersion::Mv1,
            Ok(2) => MountVersion::Mv2,
            Ok(n) => MountVersion::Other(n),
            Err(_) => MountVersion::Mv0,
        }
    }
}

/// Vault operations supported by Vaku. Used to drive path rewriting.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum VaultOp {
    List,
    Read,
    Write,
    Delete,
    Destroy,
    DeleteMeta,
    ReadMeta,
}

/// Returns true when the operation is supported on the given mount version.
pub fn mount_supports_operation(op: VaultOp, v: MountVersion) -> bool {
    if matches!(v, MountVersion::Mv2) {
        return true;
    }
    !matches!(op, VaultOp::Destroy | VaultOp::DeleteMeta | VaultOp::ReadMeta)
}

/// Find the mount that `path` lives under. Returns `(mount_path, version)`.
pub async fn mount_info(
    provider: &dyn crate::api::mount_provider::MountProvider,
    path: &str,
) -> Result<(String, MountVersion), Error> {
    let mounts = provider.list_mounts().await.map_err(|e| {
        Error::wrap(
            path,
            ErrorKind::MountInfo,
            Some(Box::new(Error::wrap(
                &e.to_string(),
                ErrorKind::ListMounts,
                None,
            ))),
        )
    })?;

    for mount in mounts {
        // Always check against the folder form so we don't match a prefix like
        // `foo/bar/` against the path `foo/barbar/...`.
        let mp = ensure_folder(&mount.path);
        if path.starts_with(&mp) {
            let version = if mount.version.is_empty() {
                MountVersion::Mv0
            } else {
                MountVersion::from_str(&mount.version)
            };
            return Ok((mp, version));
        }
    }

    Err(Error::wrap(
        path,
        ErrorKind::MountInfo,
        Some(Box::new(Error::wrap(path, ErrorKind::NoMount, None))),
    ))
}

/// Rewrite `p` for the given operation, transparently inserting KV v2's
/// `data` / `metadata` / `destroy` sub-paths when required.
pub async fn rewrite_path(
    provider: &dyn crate::api::mount_provider::MountProvider,
    p: &str,
    op: VaultOp,
) -> Result<(String, MountVersion), Error> {
    let (mount, version) = match mount_info(provider, p).await {
        Ok(v) => v,
        Err(e) => return Err(Error::wrap(p, ErrorKind::RewritePath, Some(Box::new(e)))),
    };

    if !mount_supports_operation(op, version) {
        return Err(Error::wrap(p, ErrorKind::MountVersion, None));
    }

    if !matches!(version, MountVersion::Mv2) {
        return Ok((p.to_string(), version));
    }

    let rewritten = match op {
        VaultOp::List | VaultOp::DeleteMeta | VaultOp::ReadMeta => {
            insert_into_path(p, &mount, KV2_METADATA)
        }
        VaultOp::Read | VaultOp::Write | VaultOp::Delete => insert_into_path(p, &mount, KV2_DATA),
        VaultOp::Destroy => insert_into_path(p, &mount, KV2_DESTROY),
    };

    Ok((rewritten, version))
}

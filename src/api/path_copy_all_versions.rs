//! `PathCopyAllVersions` — replicate every historical version of a KV v2 secret.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use crate::api::mounts::{mount_info, MountVersion};
use crate::api::secret::SecretVersionMeta;
use serde_json::Map;

impl Client {
    /// Copy every version (including deleted/destroyed placeholders) of `src`
    /// to `dst`. Both sides must be KV v2.
    pub async fn path_copy_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        validate_kv2(self.src().mount_provider.as_ref(), src)
            .await
            .map_err(|e| Error::wrap(src, ErrorKind::PathCopyAllVersions, Some(Box::new(e))))?;
        validate_kv2(self.dst().mount_provider.as_ref(), dst)
            .await
            .map_err(|e| Error::wrap(dst, ErrorKind::PathCopyAllVersions, Some(Box::new(e))))?;

        let meta = self
            .path_read_meta(src)
            .await
            .map_err(|e| Error::wrap(src, ErrorKind::PathCopyAllVersions, Some(Box::new(e))))?;
        let Some(meta) = meta else { return Ok(()) };
        if meta.versions.is_empty() {
            return Ok(());
        }

        // Sort by version number ascending so destination ordering is stable.
        let mut versions: Vec<i64> = meta.versions.keys().copied().collect();
        versions.sort_unstable();

        let dst_client = self.as_destination();
        for v in versions {
            let vmeta = meta.versions.get(&v).cloned().unwrap_or_default();
            let data = self.get_version_data(src, v, &vmeta).await?;
            dst_client
                .path_write(dst, Some(data))
                .await
                .map_err(|e| Error::wrap(dst, ErrorKind::PathCopyAllVersions, Some(Box::new(e))))?;
        }

        Ok(())
    }

    async fn get_version_data(
        &self,
        src: &str,
        version: i64,
        vmeta: &SecretVersionMeta,
    ) -> Result<Map<String, serde_json::Value>, Error> {
        if vmeta.deleted || vmeta.destroyed {
            return Ok(Map::new());
        }
        let data = self
            .path_read_version(src, version)
            .await
            .map_err(|e| Error::wrap(src, ErrorKind::PathCopyAllVersions, Some(Box::new(e))))?;
        Ok(data.unwrap_or_default())
    }
}

/// Returns `Ok(())` when `path` is on a KV v2 mount, otherwise [`MountVersion`]
/// or [`MountInfo`] errors.
pub(crate) async fn validate_kv2(
    provider: &dyn crate::api::mount_provider::MountProvider,
    path: &str,
) -> Result<(), Error> {
    let (_, version) = mount_info(provider, path).await?;
    if !matches!(version, MountVersion::Mv2) {
        return Err(Error::wrap("", ErrorKind::MountVersion, None));
    }
    Ok(())
}

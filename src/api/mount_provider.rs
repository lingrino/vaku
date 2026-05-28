//! Mount discovery: figure out which mount path a given path belongs to and
//! what KV version that mount uses. Required for path rewriting on KV v2.

use crate::api::error::{Error, ErrorKind};
use crate::api::logical::Logical;
use crate::api::secret::KV2_VERSION;
use async_trait::async_trait;
use serde_json::Value;
use std::sync::Arc;

/// A high-level representation of a Vault mount.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Mount {
    pub path: String,
    pub r#type: String,
    pub version: String,
}

/// Trait for objects that can list Vault mounts. The default implementation
/// uses `sys/mounts`; users without permission can plug in a
/// [`StaticMountProvider`].
#[async_trait]
pub trait MountProvider: Send + Sync {
    async fn list_mounts(&self) -> Result<Vec<Mount>, Error>;
}

/// `sys/mounts`-based mount provider. Internal type — exposed via the [`Client`]
/// builder.
#[derive(Clone)]
pub(crate) struct DefaultMountProvider {
    pub(crate) logical: Arc<dyn Logical>,
}

#[async_trait]
impl MountProvider for DefaultMountProvider {
    async fn list_mounts(&self) -> Result<Vec<Mount>, Error> {
        let secret = self.logical.read("sys/mounts").await.map_err(|e| {
            Error::wrap("", ErrorKind::MountInfo, Some(Box::new(Error::wrap(
                &e.to_string(),
                ErrorKind::ListMounts,
                None,
            ))))
        })?;

        let Some(secret) = secret else { return Ok(Vec::new()) };
        let Some(data) = secret.data else { return Ok(Vec::new()) };

        let mut mounts = Vec::with_capacity(data.len());
        // The `sys/mounts` endpoint can place mounts at the top level (legacy)
        // or under a nested `data` key (newer Vault). Walk both shapes.
        let candidates: Vec<(&String, &Value)> = if let Some(nested) = data.get("data").and_then(Value::as_object) {
            nested.iter().collect()
        } else {
            data.iter().filter(|(_, v)| v.is_object()).collect()
        };

        for (mount_path, value) in candidates {
            let obj = match value.as_object() {
                Some(o) => o,
                None => continue,
            };
            // Skip Vault keys like "request_id", "wrap_info", etc.
            if !obj.contains_key("type") {
                continue;
            }
            let r#type = obj.get("type").and_then(Value::as_str).unwrap_or_default().to_string();
            let version = obj
                .get("options")
                .and_then(Value::as_object)
                .and_then(|o| o.get(KV2_VERSION))
                .and_then(Value::as_str)
                .unwrap_or_default()
                .to_string();
            mounts.push(Mount {
                path: mount_path.clone(),
                r#type,
                version,
            });
        }

        Ok(mounts)
    }
}

/// A mount provider that returns a single, statically-configured mount.
/// Useful when the Vault token lacks permission to list mounts but the caller
/// knows the mount path + version.
#[derive(Debug, Clone)]
pub struct StaticMountProvider {
    mount: Mount,
}

impl StaticMountProvider {
    pub fn new(path: impl Into<String>, version: impl Into<String>) -> Self {
        Self {
            mount: Mount {
                path: path.into(),
                r#type: "kv".to_string(),
                version: version.into(),
            },
        }
    }
}

#[async_trait]
impl MountProvider for StaticMountProvider {
    async fn list_mounts(&self) -> Result<Vec<Mount>, Error> {
        Ok(vec![self.mount.clone()])
    }
}

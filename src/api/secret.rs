//! KV secret types and decoders.

use serde_json::Map;
use serde_json::Value;
use std::collections::BTreeMap;

/// KV v2 keys/sub-paths used in path rewriting.
pub(crate) const KV2_DATA: &str = "data";
pub(crate) const KV2_METADATA: &str = "metadata";
pub(crate) const KV2_DESTROY: &str = "destroy";
pub(crate) const KV2_VERSION: &str = "version";
pub(crate) const KV2_VERSIONS: &str = "versions";

/// Metadata for a single version of a secret.
#[derive(Debug, Clone, Default, PartialEq, Eq)]
pub struct SecretVersionMeta {
    pub created_time: String,
    pub deleted: bool,
    pub destroyed: bool,
}

/// Aggregated metadata for a secret across all of its versions.
#[derive(Debug, Clone, Default, PartialEq, Eq)]
pub struct SecretMeta {
    pub current_version: i64,
    pub versions: BTreeMap<i64, SecretVersionMeta>,
}

/// Pull the inner `data` map out of a KV v2 read response, returning `None`
/// for deleted or destroyed secrets.
pub(crate) fn extract_v2_read(data: Option<&Map<String, Value>>) -> Option<Map<String, Value>> {
    let data = data?;
    if is_deleted(data) {
        return None;
    }
    let inner = data.get(KV2_DATA)?;
    let inner_map = inner.as_object()?.clone();
    if inner_map.is_empty() && inner.is_object() {
        // Vault returns an empty `data: {}` for destroyed-by-meta and
        // never-written secrets. Mirror Go's behaviour: still surface the
        // empty map so callers can distinguish from "nil" read.
        // (Go code: returns the inner map even if empty.)
        return Some(inner_map);
    }
    Some(inner_map)
}

/// True when KV v2 metadata indicates a deleted or destroyed secret.
pub(crate) fn is_deleted(data: &Map<String, Value>) -> bool {
    let Some(metadata) = data.get(KV2_METADATA).and_then(Value::as_object) else {
        return true;
    };
    let deletion_time_ok =
        matches!(metadata.get("deletion_time"), Some(Value::String(s)) if s.is_empty());
    if !deletion_time_ok {
        return true;
    }
    let destroyed_ok = matches!(metadata.get("destroyed"), Some(Value::Bool(false)));
    if !destroyed_ok {
        return true;
    }
    false
}

/// Parse a KV v2 metadata response into a [`SecretMeta`].
pub(crate) fn extract_secret_meta(data: Option<&Map<String, Value>>) -> SecretMeta {
    let mut meta = SecretMeta::default();
    let Some(data) = data else { return meta };

    if let Some(v) = data.get("current_version") {
        meta.current_version = extract_int(v);
    }

    let Some(versions) = data.get(KV2_VERSIONS).and_then(Value::as_object) else {
        return meta;
    };

    for (version_str, raw) in versions {
        let Ok(version) = version_str.parse::<i64>() else {
            continue;
        };
        let Some(version_data) = raw.as_object() else {
            continue;
        };

        let mut vm = SecretVersionMeta::default();

        if let Some(Value::String(s)) = version_data.get("created_time") {
            vm.created_time = s.clone();
        }
        if let Some(Value::String(s)) = version_data.get("deletion_time") {
            if !s.is_empty() {
                vm.deleted = true;
            }
        }
        if let Some(Value::Bool(b)) = version_data.get("destroyed") {
            vm.destroyed = *b;
        }

        meta.versions.insert(version, vm);
    }

    meta
}

/// Coerce a JSON value into an `i64`. Vault sometimes serializes numbers as
/// floats; this mirrors Go's `extractInt`.
pub(crate) fn extract_int(v: &Value) -> i64 {
    match v {
        Value::Number(n) => {
            if let Some(i) = n.as_i64() {
                i
            } else if let Some(f) = n.as_f64() {
                f as i64
            } else {
                0
            }
        }
        _ => 0,
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json::json;

    fn obj(v: Value) -> Option<Map<String, Value>> {
        v.as_object().cloned()
    }

    #[test]
    fn extract_v2_read_nil() {
        assert!(extract_v2_read(None).is_none());
    }

    #[test]
    fn extract_v2_read_no_metadata() {
        let data = obj(json!({"foo": "bar"})).unwrap();
        assert!(extract_v2_read(Some(&data)).is_none());
    }

    #[test]
    fn extract_v2_read_meta_only() {
        let data = obj(json!({"metadata": {"foo": "bar"}})).unwrap();
        assert!(extract_v2_read(Some(&data)).is_none());
    }

    #[test]
    fn extract_v2_read_missing_destroyed() {
        let data = obj(json!({"metadata": {"deletion_time": ""}})).unwrap();
        assert!(extract_v2_read(Some(&data)).is_none());
    }

    #[test]
    fn extract_v2_read_no_data_field() {
        let data = obj(json!({
            "metadata": {"deletion_time": "", "destroyed": false}
        }))
        .unwrap();
        assert!(extract_v2_read(Some(&data)).is_none());
    }

    #[test]
    fn extract_v2_read_happy() {
        let data = obj(json!({
            "metadata": {"deletion_time": "", "destroyed": false},
            "data": {"foo": "bar"},
        }))
        .unwrap();
        let got = extract_v2_read(Some(&data)).unwrap();
        assert_eq!(got.get("foo").unwrap(), &Value::String("bar".into()));
    }

    #[test]
    fn extract_secret_meta_full() {
        let data = obj(json!({
            "current_version": 3,
            "versions": {
                "1": {"created_time": "2023-01-01T00:00:00Z", "deletion_time": "", "destroyed": false},
                "2": {"created_time": "2023-01-02T00:00:00Z", "deletion_time": "2023-01-03T00:00:00Z", "destroyed": false},
                "3": {"created_time": "2023-01-04T00:00:00Z", "deletion_time": "", "destroyed": true},
            }
        })).unwrap();
        let got = extract_secret_meta(Some(&data));
        assert_eq!(got.current_version, 3);
        assert_eq!(got.versions.len(), 3);
        assert!(!got.versions[&1].deleted);
        assert!(got.versions[&2].deleted);
        assert!(got.versions[&3].destroyed);
    }

    #[test]
    fn extract_secret_meta_nil_data() {
        let got = extract_secret_meta(None);
        assert_eq!(got.current_version, 0);
        assert!(got.versions.is_empty());
    }

    #[test]
    fn extract_secret_meta_empty_versions() {
        let data = obj(json!({"current_version": 0, "versions": {}})).unwrap();
        let got = extract_secret_meta(Some(&data));
        assert_eq!(got.current_version, 0);
        assert!(got.versions.is_empty());
    }
}

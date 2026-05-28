//! Path-marker-driven `Logical` injector used by tests to splice canned
//! responses and errors into the Vault HTTP surface — a direct port of Go's
//! `logicalInjector`.
//!
//! How it works: callers embed an injection trail into the path itself:
//!
//!   `something/<name>/<op>/inject/some/more`
//!
//! When the operation matches (`read`, `write`, `list`, `delete`), the
//! injector returns the canned secret/error registered under `<name>`. When
//! the op doesn't match, it strips the marker and forwards to the real
//! Logical. For LIST the marker is re-embedded into returned child keys so
//! deeper recursive operations keep firing.

use crate::common::seeds::shared_clients;
use async_trait::async_trait;
use serde_json::{Map, Value};
use std::error::Error as StdError;
use std::fmt;
use std::sync::Arc;
use vaku::api::error::Error as VakuError;
use vaku::api::helpers::{ensure_folder, is_folder, path_join};
use vaku::api::logical::{Logical, Secret};

/// One injection slot.
#[derive(Clone, Default)]
pub struct Inject {
    pub secret: Option<Secret>,
    pub err: Option<&'static str>,
}

/// Static error used by injectors named "error". Matches Go's `errInject`.
pub const ERR_INJECT_MSG: &str = "injected error";

#[derive(Debug)]
pub struct InjectedErr(pub &'static str);
impl fmt::Display for InjectedErr {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.0)
    }
}
impl StdError for InjectedErr {}

fn injects(name: &str) -> Option<Inject> {
    let mut data = Map::new();
    match name {
        "error" => Some(Inject {
            secret: None,
            err: Some(ERR_INJECT_MSG),
        }),
        "nildata" => Some(Inject {
            secret: Some(Secret { data: None }),
            err: None,
        }),
        "nilkeys" => {
            data.insert("keys".into(), Value::Null);
            Some(Inject {
                secret: Some(Secret { data: Some(data) }),
                err: None,
            })
        }
        "intkeys" => {
            data.insert("keys".into(), Value::from(1));
            Some(Inject {
                secret: Some(Secret { data: Some(data) }),
                err: None,
            })
        }
        "listintkeys" => {
            data.insert("keys".into(), Value::Array(vec![Value::from(1)]));
            Some(Inject {
                secret: Some(Secret { data: Some(data) }),
                err: None,
            })
        }
        "funcdata" => {
            // Used by the search tests to make json marshal fail. In Go the
            // value is a `func(){}` which `encoding/json` can't marshal; for
            // Rust we represent unmarshalable data by inserting a number that
            // makes the JSON valid but the *value* contains a non-finite f64
            // (NaN), which serde_json::to_string refuses.
            let mut inner = Map::new();
            inner.insert(
                "foo".into(),
                Value::Number(
                    serde_json::Number::from_f64(f64::NAN)
                        .unwrap_or_else(|| serde_json::Number::from(0)),
                ),
            );
            let mut outer = Map::new();
            outer.insert("data".into(), Value::Object(inner));
            let mut meta = Map::new();
            meta.insert("destroyed".into(), Value::Bool(false));
            meta.insert("deletion_time".into(), Value::String(String::new()));
            outer.insert("metadata".into(), Value::Object(meta));
            Some(Inject {
                secret: Some(Secret { data: Some(outer) }),
                err: None,
            })
        }
        _ => None,
    }
}

/// Wraps another `Logical` and rewrites paths / returns canned responses
/// based on embedded markers.
pub struct LogicalInjector {
    inner: Arc<dyn Logical>,
    disabled: bool,
}

impl LogicalInjector {
    pub fn new(inner: Arc<dyn Logical>, disabled: bool) -> Self {
        Self { inner, disabled }
    }

    /// Resolve a path against the injection rules. Returns `(clean_path, hit?)`.
    /// When `hit` is `Some`, the caller should return the canned response.
    fn run(&self, p: &str, op: &str) -> (String, Option<Inject>) {
        // Strip trailing slash for matching but remember it.
        let trimmed = p.trim_end_matches('/');

        let parts: Vec<&str> = trimmed.split('/').collect();
        let mut inject_op: Option<String> = None;
        let mut inject_name: Option<String> = None;
        let mut cleaned: Vec<String> = parts.iter().map(|s| (*s).to_string()).collect();
        for i in 0..parts.len() {
            if parts[i] == "inject" && i >= 2 {
                inject_op = Some(parts[i - 1].to_string());
                inject_name = Some(parts[i - 2].to_string());
                cleaned[i] = String::new();
                cleaned[i - 1] = String::new();
                cleaned[i - 2] = String::new();
            }
        }

        let Some(name) = inject_name else {
            return (p.to_string(), None);
        };

        // path_join with empty parts collapses them.
        let collapsed = path_join(&cleaned.iter().map(String::as_str).collect::<Vec<_>>());
        let clean_path = if is_folder(p) {
            ensure_folder(&collapsed)
        } else {
            collapsed
        };

        let want_op = inject_op.unwrap_or_default();
        if want_op != op {
            return (clean_path, None);
        }

        if self.disabled {
            return (clean_path, None);
        }

        match injects(&name) {
            Some(inj) => (clean_path, Some(inj)),
            None => (clean_path, None),
        }
    }

    /// LIST-specific result fixup: re-add the injection markers to listed keys
    /// so they continue to fire on subsequent operations down the tree.
    fn rewrite_list_result(orig_path: &str, sec: Option<Secret>) -> Option<Secret> {
        let mut sec = sec?;
        let Some(data) = &sec.data else {
            return Some(sec);
        };
        let Some(Value::Array(arr)) = data.get("keys") else {
            return Some(sec);
        };

        // Only fire when the original path ended in `.../inject` (or
        // `.../inject/`).
        let trimmed = orig_path.trim_end_matches('/');
        let base = std::path::Path::new(trimmed)
            .file_name()
            .and_then(|s| s.to_str())
            .unwrap_or_default();
        if base != "inject" {
            return Some(sec);
        }

        // Re-derive the `name/op/inject` triple from the original path.
        let parts: Vec<&str> = trimmed.rsplitn(4, '/').collect();
        // parts is [inject, op, name, rest...] reversed.
        if parts.len() < 3 {
            return Some(sec);
        }
        let triple = path_join(&[parts[2], parts[1], parts[0]]);

        let mut new_keys = Vec::with_capacity(arr.len());
        for v in arr {
            let Value::String(key) = v else {
                new_keys.push(v.clone());
                continue;
            };
            let trailing = key.ends_with('/');
            let stripped: &str = key.strip_suffix('/').unwrap_or(key);
            let mut new_key = path_join(&[stripped, &triple]);
            if trailing && !new_key.ends_with('/') {
                new_key.push('/');
            }
            new_keys.push(Value::String(new_key));
        }
        let mut new_data = data.clone();
        new_data.insert("keys".into(), Value::Array(new_keys));
        sec.data = Some(new_data);
        Some(sec)
    }
}

#[async_trait]
impl Logical for LogicalInjector {
    async fn read(&self, p: &str) -> Result<Option<Secret>, vaku::api::error::BoxError> {
        let (clean, inj) = self.run(p, "read");
        if let Some(inj) = inj {
            return finalize(inj);
        }
        self.inner.read(&clean).await
    }

    async fn read_with_data(
        &self,
        p: &str,
        params: &[(&str, &str)],
    ) -> Result<Option<Secret>, vaku::api::error::BoxError> {
        let (clean, inj) = self.run(p, "read");
        if let Some(inj) = inj {
            return finalize(inj);
        }
        self.inner.read_with_data(&clean, params).await
    }

    async fn list(&self, p: &str) -> Result<Option<Secret>, vaku::api::error::BoxError> {
        let (clean, inj) = self.run(p, "list");
        if let Some(inj) = inj {
            return finalize(inj);
        }
        let sec = self.inner.list(&clean).await?;
        Ok(LogicalInjector::rewrite_list_result(p, sec))
    }

    async fn write(
        &self,
        p: &str,
        data: Value,
    ) -> Result<Option<Secret>, vaku::api::error::BoxError> {
        let (clean, inj) = self.run(p, "write");
        if let Some(inj) = inj {
            return finalize(inj);
        }
        self.inner.write(&clean, data).await
    }

    async fn delete(&self, p: &str) -> Result<Option<Secret>, vaku::api::error::BoxError> {
        let (clean, inj) = self.run(p, "delete");
        if let Some(inj) = inj {
            return finalize(inj);
        }
        self.inner.delete(&clean).await
    }
}

fn finalize(inj: Inject) -> Result<Option<Secret>, vaku::api::error::BoxError> {
    if let Some(msg) = inj.err {
        return Err(Box::new(InjectedErr(msg)));
    }
    Ok(inj.secret)
}

// Silence unused import warning when this file is reached only via tests.
#[allow(dead_code)]
fn _touch_shared() {
    let _ = shared_clients;
}

// Used by some tests to wrap the underlying error message into a VakuError so
// chain comparisons line up.
#[allow(dead_code)]
pub fn inject_error_node() -> VakuError {
    VakuError::wrap(
        ERR_INJECT_MSG,
        vaku::api::error::ErrorKind::Custom(ERR_INJECT_MSG.into()),
        None,
    )
}

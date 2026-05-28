//! Path-manipulation helpers shared across the API.
//!
//! Semantics mirror the Go `api/helpers.go`. In particular, [`path_join`]
//! collapses repeated slashes and trims a leading slash but preserves a single
//! trailing slash when the **last** input had one.

use serde_json::Map;
use serde_json::Value;
use std::collections::BTreeMap;

/// A 1:1 port of Go's `path.Clean` that keeps the behaviour Vaku relies on:
/// collapse `//` and `.` segments, strip leading `/`. It does **not** preserve
/// trailing slashes — that's [`path_join`]'s job.
fn go_path_clean(p: &str) -> String {
    if p.is_empty() {
        return ".".to_string();
    }
    let rooted = p.starts_with('/');

    // Lazy multi-pass cleanup using a stack-of-segment approach.
    let mut out: Vec<&str> = Vec::new();
    for seg in p.split('/') {
        match seg {
            "" | "." => {}
            ".." => {
                // Go's path.Clean drops ".." against the root, but otherwise
                // pops the previous segment.
                if let Some(last) = out.last() {
                    if *last != ".." {
                        out.pop();
                        continue;
                    }
                }
                if !rooted {
                    out.push("..");
                }
            }
            s => out.push(s),
        }
    }

    let mut joined = out.join("/");
    if rooted {
        joined.insert(0, '/');
    }
    if joined.is_empty() {
        ".".to_string()
    } else {
        joined
    }
}

/// Combine multiple path segments into one Vaku-style path:
/// collapse repeated slashes, drop a leading slash, but preserve a single
/// trailing slash when the **last** segment ends with `/`.
pub fn path_join(parts: &[&str]) -> String {
    if parts.is_empty() {
        return String::new();
    }
    let last_has_trailing = parts.last().is_some_and(|p| p.ends_with('/'));

    // Skip purely-empty inputs so behaviour matches Go's path.Join, which
    // ignores them.
    let filtered: Vec<&str> = parts.iter().copied().filter(|s| !s.is_empty()).collect();
    if filtered.is_empty() {
        return String::new();
    }

    let cleaned = go_path_clean(&filtered.join("/"));
    let mut s: String = cleaned.trim_start_matches('/').to_string();

    if last_has_trailing && !s.ends_with('/') {
        s.push('/');
    }
    s
}

/// True when `p` is a folder (i.e. ends with `/`).
pub fn is_folder(p: &str) -> bool {
    p.ends_with('/')
}

/// Ensure `p` is a folder by adding a trailing `/` if missing.
pub fn ensure_folder(p: &str) -> String {
    path_join(&[p, "/"])
}

/// Add `prefix` to the start of `p`.
pub fn add_prefix(p: &str, prefix: &str) -> String {
    path_join(&[prefix, p])
}

/// Add `prefix` to `p` only if it isn't already prefixed.
pub fn ensure_prefix(p: &str, prefix: &str) -> String {
    if p.starts_with(prefix) {
        p.to_string()
    } else {
        path_join(&[prefix, p])
    }
}

/// Add `prefix` to every entry in the list.
pub fn add_prefix_list(list: &mut [String], prefix: &str) {
    for item in list.iter_mut() {
        *item = path_join(&[prefix, item]);
    }
}

/// Add `prefix` to every entry that doesn't already have it.
pub fn ensure_prefix_list(list: &mut [String], prefix: &str) {
    for item in list.iter_mut() {
        if !item.starts_with(prefix) {
            *item = path_join(&[prefix, item]);
        }
    }
}

/// Remove `prefix` from the start of every entry. Entries without the prefix
/// are returned cleaned (matching Go's `PathJoin(strings.TrimPrefix(...))`).
pub fn trim_prefix_list(list: &mut [String], prefix: &str) {
    for item in list.iter_mut() {
        let trimmed: &str = item.strip_prefix(prefix).unwrap_or(item.as_str());
        *item = path_join(&[trimmed]);
    }
}

/// Ensure every map key has `prefix`.
pub fn ensure_prefix_map(map: &mut BTreeMap<String, Map<String, Value>>, prefix: &str) {
    let keys: Vec<String> = map.keys().cloned().collect();
    for k in keys {
        let new_key = ensure_prefix(&k, prefix);
        if new_key != k {
            if let Some(v) = map.remove(&k) {
                map.insert(new_key, v);
            }
        }
    }
}

/// Trim `prefix` from every map key. Keys without the prefix are cleaned via
/// `path_join` (mirrors Go's behaviour).
pub fn trim_prefix_map(map: &mut BTreeMap<String, Map<String, Value>>, prefix: &str) {
    let keys: Vec<String> = map.keys().cloned().collect();
    for k in keys {
        let trimmed: String = match k.strip_prefix(prefix) {
            Some(s) => path_join(&[s]),
            None => path_join(&[&k]),
        };
        if trimmed != k {
            if let Some(v) = map.remove(&k) {
                map.insert(trimmed, v);
            }
        }
    }
}

/// Insert `insert` into `path` immediately after `after`.
pub fn insert_into_path(path: &str, after: &str, insert: &str) -> String {
    let tail: &str = path.strip_prefix(after).unwrap_or(path);
    path_join(&[after, insert, tail])
}

/// Merge `b` into `a`, preferring entries from `b`.
pub(crate) fn merge_maps(
    a: &mut BTreeMap<String, Map<String, Value>>,
    b: BTreeMap<String, Map<String, Value>>,
) {
    for (k, v) in b {
        a.insert(k, v);
    }
}

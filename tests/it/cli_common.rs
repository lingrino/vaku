//! Shared CLI testing helpers — a mock [`ClientInterface`] that returns
//! canned values identical to Go's `testVakuClient`.

use async_trait::async_trait;
use serde_json::{json, Map, Value};
use std::collections::BTreeMap;
use std::sync::Arc;
use vaku::api::error::Error;
use vaku::api::secret::{SecretMeta, SecretVersionMeta};
use vaku::cli::client_iface::ClientInterface;
use vaku::cli::runner;

pub fn run_cli(args: &[&str]) -> (i32, String, String) {
    run_cli_with_client(args, None)
}

pub fn run_cli_with_mock(args: &[&str]) -> (i32, String, String) {
    let mock: Arc<dyn ClientInterface> = Arc::new(MockClient {});
    run_cli_with_client(args, Some(mock))
}

pub fn run_cli_with_client(
    args: &[&str],
    client: Option<Arc<dyn ClientInterface>>,
) -> (i32, String, String) {
    let mut out: Vec<u8> = Vec::new();
    let mut err: Vec<u8> = Vec::new();
    // Mirror Go's cmd test harness which sets `cli.flagIndent = ""`. Override
    // here unless the caller passed --indent-char/-i themselves.
    let mut argv: Vec<String> = Vec::with_capacity(args.len() + 2);
    let already_set = args.iter().any(|a| {
        *a == "-i"
            || *a == "--indent-char"
            || a.starts_with("-i")
            || a.starts_with("--indent-char=")
    });
    if !already_set {
        argv.push("--indent-char".into());
        argv.push("".into());
    }
    argv.extend(args.iter().map(|s| s.to_string()));
    let code = runner::run_with_client("dev", &argv, &mut out, &mut err, client);
    (
        code as i32,
        String::from_utf8(out).unwrap_or_default(),
        String::from_utf8(err).unwrap_or_default(),
    )
}

pub struct MockClient;

fn m(kvs: &[(&str, Value)]) -> Map<String, Value> {
    let mut m = Map::new();
    for (k, v) in kvs {
        m.insert((*k).to_string(), v.clone());
    }
    m
}

#[async_trait]
impl ClientInterface for MockClient {
    async fn path_list(&self, _: &str) -> Result<Vec<String>, Error> {
        Ok(vec!["foo".into(), "moo".into()])
    }
    async fn path_read(&self, _: &str) -> Result<Option<Map<String, Value>>, Error> {
        Ok(Some(m(&[("biz", json!("baz")), ("foo", json!("bar"))])))
    }
    async fn path_read_meta(&self, _: &str) -> Result<Option<SecretMeta>, Error> {
        let mut versions = BTreeMap::new();
        versions.insert(1, SecretVersionMeta::default());
        Ok(Some(SecretMeta {
            current_version: 1,
            versions,
        }))
    }
    async fn path_read_version(
        &self,
        _: &str,
        _: i64,
    ) -> Result<Option<Map<String, Value>>, Error> {
        Ok(Some(m(&[("biz", json!("baz")), ("foo", json!("bar"))])))
    }
    async fn path_write(&self, _: &str, _: Option<Map<String, Value>>) -> Result<(), Error> {
        Ok(())
    }
    async fn path_delete(&self, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn path_delete_meta(&self, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn path_destroy(&self, _: &str, _: &[i64]) -> Result<(), Error> {
        Ok(())
    }
    async fn path_update(&self, _: &str, _: Option<Map<String, Value>>) -> Result<(), Error> {
        Ok(())
    }
    async fn path_search(&self, _: &str, _: &str) -> Result<bool, Error> {
        Ok(true)
    }
    async fn path_copy(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn path_copy_all_versions(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn path_move(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn path_move_all_versions(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }

    async fn folder_list(&self, _: &str) -> Result<Vec<String>, Error> {
        Ok(vec!["foo/bar".into(), "foo/baz".into(), "bim/bom".into()])
    }
    async fn folder_read(
        &self,
        _: &str,
    ) -> Result<Option<BTreeMap<String, Map<String, Value>>>, Error> {
        Ok(Some(BTreeMap::from([
            (
                "foo".to_string(),
                m(&[("bim", json!("bom")), ("biz", json!("baz"))]),
            ),
            ("bar".to_string(), m(&[("hoo", json!("boo"))])),
        ])))
    }
    async fn folder_write(&self, _: BTreeMap<String, Map<String, Value>>) -> Result<(), Error> {
        Ok(())
    }
    async fn folder_delete(&self, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn folder_delete_meta(&self, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn folder_destroy(&self, _: &str, _: &[i64]) -> Result<(), Error> {
        Ok(())
    }
    async fn folder_search(&self, _: &str, _: &str) -> Result<Vec<String>, Error> {
        Ok(vec!["foo/bar".into(), "bim/bom".into()])
    }
    async fn folder_copy(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn folder_copy_all_versions(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn folder_move(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }
    async fn folder_move_all_versions(&self, _: &str, _: &str) -> Result<(), Error> {
        Ok(())
    }
}

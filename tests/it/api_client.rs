//! Ports `api/client_test.go` — pure-function tests for the builder and
//! the input/output path helpers.

use vaku::api::client::Client;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};

#[tokio::test(flavor = "multi_thread")]
async fn test_new_client_defaults() {
    let c = Client::builder().build().unwrap();
    assert_eq!(c.workers(), 10);
    assert!(!c.absolute_path());
    assert!(!c.ignore_access_errors());
}

#[tokio::test(flavor = "multi_thread")]
async fn test_new_client_bad_workers() {
    let err = Client::builder().with_workers(0).build().map(|_| ()).unwrap_err();
    let de: &(dyn std::error::Error + 'static) = &err;
    compare_errors(Some(de), &[
        ErrorKind::ApplyOptions.into(),
        ErrorKind::NumWorkers.into(),
    ]);
}

#[tokio::test(flavor = "multi_thread")]
async fn test_new_client_options() {
    let c = Client::builder()
        .with_workers(100)
        .with_absolute_path(true)
        .with_ignore_access_errors(true)
        .build()
        .unwrap();
    assert_eq!(c.workers(), 100);
    assert!(c.absolute_path());
    assert!(c.ignore_access_errors());
}

#[tokio::test(flavor = "multi_thread")]
async fn test_input_output_paths() {
    use serde_json::{json, Map, Value};
    use std::collections::BTreeMap;

    // input_path: absolute returns as-is; non-absolute prepends root.
    let abs = Client::builder().with_absolute_path(true).build().unwrap();
    let rel = Client::builder().build().unwrap();

    let abs_ip = abs.input_path("3/4", "0/1/2");
    let rel_ip = rel.input_path("3/4", "0/1/2");
    assert_eq!(abs_ip, "3/4");
    assert_eq!(rel_ip, "0/1/2/3/4");

    // output_path: absolute ensures the prefix; non-absolute trims it.
    let abs_op = abs.output_path("3", "0/1/2");
    let rel_op = rel.output_path("3", "0/1/2");
    assert_eq!(abs_op, "0/1/2/3");
    assert_eq!(rel_op, "3");

    // output_paths
    let mut paths = vec!["3".to_string(), "4".to_string()];
    abs.output_paths(&mut paths, "0/1/2");
    assert_eq!(paths, vec!["0/1/2/3".to_string(), "0/1/2/4".to_string()]);

    let mut paths = vec!["3".to_string(), "4".to_string()];
    rel.output_paths(&mut paths, "0/1/2");
    assert_eq!(paths, vec!["3".to_string(), "4".to_string()]);

    // swap_paths
    let mut data: BTreeMap<String, Map<String, Value>> = BTreeMap::from([
        ("0/1/2/3".to_string(), Map::from_iter([("k".to_string(), json!("v"))])),
        ("0/1/2/4".to_string(), Map::from_iter([("k".to_string(), json!("v"))])),
    ]);
    abs.swap_paths(&mut data, "0/1/2", "00/01/02");
    assert!(data.contains_key("00/01/02/3"));
    assert!(data.contains_key("00/01/02/4"));

    let mut data: BTreeMap<String, Map<String, Value>> = BTreeMap::from([
        ("0/1/2/3".to_string(), Map::from_iter([("k".to_string(), json!("v"))])),
        ("0/1/2/4".to_string(), Map::from_iter([("k".to_string(), json!("v"))])),
    ]);
    rel.swap_paths(&mut data, "0/1/2", "00/01/02");
    assert!(data.contains_key("00/01/02/0/1/2/3"));
    assert!(data.contains_key("00/01/02/0/1/2/4"));
}

#[allow(dead_code)]
fn _err_matchers() -> Vec<ErrMatch> { Vec::new() } // touch ErrMatch

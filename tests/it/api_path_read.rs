//! Ports `api/path_read_test.go`.

use crate::common::{seeded_prefixes, shared_clients, MOUNTLESS};
use crate::skip_if_no_docker;
use serde_json::{json, Map, Value};
use std::sync::Arc;
use vaku::api::client::Client;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;
use vaku::api::logical::VaultHttpClient;

fn want_map(kvs: &[(&str, &str)]) -> Option<Map<String, Value>> {
    let mut m = Map::new();
    for (k, v) in kvs {
        m.insert((*k).to_string(), json!(*v));
    }
    Some(m)
}

struct Case {
    give: &'static str,
    want: Option<Map<String, Value>>,
    want_err: Vec<ErrMatch>,
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_read() {
    skip_if_no_docker!();
    let cases = vec![
        Case {
            give: "0/1",
            want: want_map(&[("2", "3")]),
            want_err: vec![],
        },
        Case {
            give: "0/4/13/17",
            want: want_map(&[("18", "19"), ("20", "21"), ("22", "23")]),
            want_err: vec![],
        },
        Case {
            give: "fake",
            want: None,
            want_err: vec![],
        },
        Case {
            give: MOUNTLESS,
            want: None,
            want_err: vec![
                ErrorKind::PathRead.into(),
                ErrorKind::RewritePath.into(),
                ErrorKind::MountInfo.into(),
                ErrorKind::NoMount.into(),
            ],
        },
        Case {
            give: "error/read/inject",
            want: None,
            want_err: vec![ErrorKind::PathRead.into(), ErrorKind::VaultRead.into()],
        },
    ];

    let clients = shared_clients().await;
    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let result = clients.vaku.path_read(&p).await;
            let got = match &result {
                Ok(v) => v.clone(),
                Err(_) => None,
            };
            assert_eq!(got, tt.want, "give={} prefix={}", tt.give, prefix);

            let err_ref: Option<&(dyn std::error::Error + 'static)> = match result.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_read_ignore_errors() {
    skip_if_no_docker!();
    let prefix = seeded_prefixes("error/read/inject")
        .await
        .into_iter()
        .next()
        .unwrap();
    let server = &crate::common::seeds::SERVERS.src;
    let http = Arc::new(VaultHttpClient::new(&server.addr, &server.token, None).unwrap());
    let injector = Arc::new(crate::common::injector::LogicalInjector::new(http, false));
    let client = Client::builder()
        .with_logical(injector)
        .with_ignore_access_errors(true)
        .build()
        .unwrap();

    let p = path_join(&[&prefix, "error/read/inject"]);
    let got = client.path_read(&p).await.expect("should ignore error");
    assert!(got.is_none());
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_read_version() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("readversion").await {
        if prefix.starts_with("kv1/") {
            let p = path_join(&[&prefix, "0/1"]);
            let err = clients.vaku.path_read_version(&p, 1).await.unwrap_err();
            let dyn_err: &(dyn std::error::Error + 'static) = &err;
            compare_errors(
                Some(dyn_err),
                &[
                    ErrorKind::PathReadVersion.into(),
                    ErrorKind::MountVersion.into(),
                ],
            );
            continue;
        }
        // KV2 cases
        let src_path = path_join(&[&prefix, "readversion/basic"]);
        clients
            .clean
            .path_write(&src_path, want_map(&[("version", "1")]))
            .await
            .unwrap();
        clients
            .clean
            .path_write(&src_path, want_map(&[("version", "2")]))
            .await
            .unwrap();

        let v1 = clients.vaku.path_read_version(&src_path, 1).await.unwrap();
        assert_eq!(v1, want_map(&[("version", "1")]));
        let v2 = clients.vaku.path_read_version(&src_path, 2).await.unwrap();
        assert_eq!(v2, want_map(&[("version", "2")]));

        let non = clients
            .vaku
            .path_read_version(&src_path, 999)
            .await
            .unwrap();
        assert!(non.is_none());

        let err = clients
            .vaku
            .path_read_version(&path_join(&[&prefix, "error/read/inject"]), 1)
            .await
            .unwrap_err();
        let dyn_err: &(dyn std::error::Error + 'static) = &err;
        compare_errors(
            Some(dyn_err),
            &[
                ErrorKind::PathReadVersion.into(),
                ErrorKind::VaultRead.into(),
            ],
        );
    }
}

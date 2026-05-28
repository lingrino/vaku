//! Ports `api/folder_read_test.go`.

use crate::common::{seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use serde_json::{json, Map, Value};
use std::collections::BTreeMap;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_map};

fn inner(kvs: &[(&str, &str)]) -> Map<String, Value> {
    let mut m = Map::new();
    for (k, v) in kvs {
        m.insert((*k).to_string(), json!(*v));
    }
    m
}

fn mk(items: &[(&str, &[(&str, &str)])]) -> BTreeMap<String, Map<String, Value>> {
    let mut m = BTreeMap::new();
    for (p, kvs) in items {
        m.insert((*p).to_string(), inner(kvs));
    }
    m
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_read() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        want: Option<BTreeMap<String, Map<String, Value>>>,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case {
            give: "0/1",
            want: None,
            want_err: vec![],
        },
        Case {
            give: "0/4/13/24/25",
            want: Some(mk(&[("26/27", &[("28", "29")])])),
            want_err: vec![],
        },
        Case {
            give: "0/4/13",
            want: Some(mk(&[
                ("14", &[("15", "16")]),
                ("17", &[("18", "19"), ("20", "21"), ("22", "23")]),
                ("24/25/26/27", &[("28", "29")]),
            ])),
            want_err: vec![],
        },
        Case {
            give: "error/list/inject",
            want: None,
            want_err: vec![
                ErrorKind::FolderRead.into(),
                ErrorKind::FolderReadChan.into(),
                ErrorKind::FolderListChan.into(),
                ErrorKind::PathList.into(),
                ErrorKind::VaultList.into(),
            ],
        },
        Case {
            give: "0/4/13/24/25/26/error/read/inject",
            want: None,
            want_err: vec![
                ErrorKind::FolderRead.into(),
                ErrorKind::FolderReadChan.into(),
                ErrorKind::PathRead.into(),
                ErrorKind::VaultRead.into(),
            ],
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.folder_read(&p).await;

            let mut got = res.as_ref().ok().cloned().flatten().unwrap_or_default();
            trim_prefix_map(&mut got, &prefix);
            let want = tt.want.clone().unwrap_or_default();
            assert_eq!(got, want, "give={} prefix={}", tt.give, prefix);

            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);
        }
    }
}

//! Ports `api/folder_delete_test.go`.

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

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_delete() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        read_back: BTreeMap<String, Map<String, Value>>,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case {
            give: "0/1",
            read_back: BTreeMap::new(),
            want_err: vec![],
        },
        Case {
            give: "0/4/13",
            read_back: BTreeMap::new(),
            want_err: vec![],
        },
        Case {
            give: "empty/path",
            read_back: BTreeMap::new(),
            want_err: vec![],
        },
        Case {
            give: "0/4/13/24/25/error/list/inject",
            read_back: BTreeMap::from_iter([("26/27".to_string(), inner(&[("28", "29")]))]),
            want_err: vec![
                ErrorKind::FolderDelete.into(),
                ErrorKind::FolderListChan.into(),
                ErrorKind::PathList.into(),
                ErrorKind::VaultList.into(),
            ],
        },
        Case {
            give: "0/4/13/24/25/26/error/delete/inject",
            read_back: BTreeMap::from_iter([("27".to_string(), inner(&[("28", "29")]))]),
            want_err: vec![
                ErrorKind::FolderDelete.into(),
                ErrorKind::PathDelete.into(),
                ErrorKind::VaultDelete.into(),
            ],
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.folder_delete(&p).await;
            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);

            let mut got = clients
                .clean
                .folder_read(&p)
                .await
                .unwrap()
                .unwrap_or_default();
            trim_prefix_map(&mut got, &prefix);
            assert_eq!(got, tt.read_back, "give={} prefix={}", tt.give, prefix);
        }
    }
}

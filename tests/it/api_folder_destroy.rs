//! Ports `api/folder_destroy_test.go`.

use crate::common::{seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use serde_json::{json, Map, Value};
use std::collections::BTreeMap;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_map};

fn inner(kvs: &[(&str, &str)]) -> Map<String, Value> {
    let mut m = Map::new();
    for (k, v) in kvs { m.insert((*k).to_string(), json!(*v)); }
    m
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_destroy() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        versions: Vec<i64>,
        read_back: BTreeMap<String, Map<String, Value>>,
        want_kv1: Vec<ErrMatch>,
        want_kv2: Vec<ErrMatch>,
    }
    let cases = vec![
        Case {
            give: "0", versions: vec![1,2,3],
            read_back: BTreeMap::new(),
            want_kv1: vec![
                ErrorKind::FolderDestroy.into(),
                ErrorKind::PathDestroy.into(),
                ErrorKind::MountVersion.into(),
            ],
            want_kv2: vec![],
        },
        Case {
            give: "0/1", versions: vec![3],
            read_back: BTreeMap::new(),
            want_kv1: vec![], want_kv2: vec![],
        },
        Case {
            give: "0/4/13/24/25/error/list/inject", versions: vec![1, 2],
            read_back: BTreeMap::from_iter([("26/27".to_string(), inner(&[("28", "29")]))]),
            want_kv1: vec![
                ErrorKind::FolderDestroy.into(),
                ErrorKind::FolderListChan.into(),
                ErrorKind::PathList.into(),
                ErrorKind::VaultList.into(),
            ],
            want_kv2: vec![
                ErrorKind::FolderDestroy.into(),
                ErrorKind::FolderListChan.into(),
                ErrorKind::PathList.into(),
                ErrorKind::VaultList.into(),
            ],
        },
        Case {
            give: "0/4/13/24/25/26/error/write/inject", versions: vec![1, 2],
            read_back: BTreeMap::from_iter([("27".to_string(), inner(&[("28", "29")]))]),
            want_kv1: vec![
                ErrorKind::FolderDestroy.into(),
                ErrorKind::PathDestroy.into(),
                ErrorKind::MountVersion.into(),
            ],
            want_kv2: vec![
                ErrorKind::FolderDestroy.into(),
                ErrorKind::PathDestroy.into(),
                ErrorKind::VaultWrite.into(),
            ],
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.folder_destroy(&p, &tt.versions).await;
            let want = if prefix.starts_with("kv1/") { &tt.want_kv1 } else { &tt.want_kv2 };
            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None, Err(e) => Some(e),
            };
            compare_errors(er, want);

            if prefix.starts_with("kv2/") {
                let mut got = clients.clean.folder_read(&p).await.unwrap().unwrap_or_default();
                trim_prefix_map(&mut got, &prefix);
                assert_eq!(got, tt.read_back, "give={}", tt.give);
            }
        }
    }
}

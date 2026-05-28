//! Ports `api/path_destroy_test.go`.

use crate::common::{seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use serde_json::json;
use serde_json::Map;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_destroy() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        versions: Vec<i64>,
        want_err: Vec<ErrMatch>,
        nil_read: bool,
    }
    let cases = vec![
        Case {
            give: "0/1",
            versions: vec![],
            want_err: vec![ErrorKind::PathDestroy.into()],
            nil_read: false,
        },
        Case {
            give: "0/1",
            versions: vec![1],
            want_err: vec![],
            nil_read: false,
        },
        Case {
            give: "0/1",
            versions: vec![2],
            want_err: vec![],
            nil_read: true,
        },
        Case {
            give: "fake",
            versions: vec![1, 2, 3, 4, 5, 6, 7],
            want_err: vec![],
            nil_read: true,
        },
        Case {
            give: "error/write/inject",
            versions: vec![1],
            want_err: vec![ErrorKind::PathDestroy.into(), ErrorKind::VaultWrite.into()],
            nil_read: false,
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            if prefix.starts_with("kv1/") {
                let err = clients.vaku.path_destroy(&p, &[1]).await.unwrap_err();
                let dyn_err: &(dyn std::error::Error + 'static) = &err;
                compare_errors(
                    Some(dyn_err),
                    &[
                        ErrorKind::PathDestroy.into(),
                        ErrorKind::MountVersion.into(),
                    ],
                );
                continue;
            }
            // KV2: overwrite first to create a new version
            let mut overwrite = Map::new();
            overwrite.insert("foo".into(), json!("bar"));
            clients
                .clean
                .path_write(&p, Some(overwrite.clone()))
                .await
                .unwrap();

            let res = clients.vaku.path_destroy(&p, &tt.versions).await;
            let err_ref: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);

            let read = clients.clean.path_read(&p).await.unwrap();
            if tt.nil_read {
                assert!(read.is_none(), "expected nil for {}", tt.give);
            } else {
                assert_eq!(read, Some(overwrite));
            }
        }
    }
}

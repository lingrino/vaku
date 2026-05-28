//! Ports `api/path_delete_meta_test.go`.

use crate::common::{seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_delete_meta() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        want_err: Vec<ErrMatch>,
        no_readback: bool,
    }
    let cases = vec![
        Case {
            give: "0/1",
            want_err: vec![],
            no_readback: false,
        },
        Case {
            give: "fake",
            want_err: vec![],
            no_readback: false,
        },
        Case {
            give: "error/delete/inject",
            want_err: vec![
                ErrorKind::PathDeleteMeta.into(),
                ErrorKind::VaultDelete.into(),
            ],
            no_readback: true,
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.path_delete_meta(&p).await;
            if prefix.starts_with("kv1/") {
                let err = res.unwrap_err();
                let dyn_err: &(dyn std::error::Error + 'static) = &err;
                compare_errors(
                    Some(dyn_err),
                    &[
                        ErrorKind::PathDeleteMeta.into(),
                        ErrorKind::MountVersion.into(),
                    ],
                );
                continue;
            }
            // KV2
            let err_ref: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);
            if !tt.no_readback {
                let r = clients.clean.path_read(&p).await.unwrap();
                assert!(r.is_none(), "expected nil read for {}", tt.give);
            }
        }
    }
}

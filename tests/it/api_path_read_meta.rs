//! Ports `api/path_read_meta_test.go`.

use crate::common::{seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_read_meta() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        want_err: Vec<ErrMatch>,
        want_nil: bool,
    }
    let cases = vec![
        Case { give: "0/1", want_err: vec![], want_nil: false },
        Case { give: "fake/path", want_err: vec![], want_nil: true },
        Case {
            give: "error/read/inject",
            want_err: vec![ErrorKind::PathReadMeta.into(), ErrorKind::VaultRead.into()],
            want_nil: false,
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            if prefix.starts_with("kv1/") {
                let err = clients.vaku.path_read_meta(&p).await.unwrap_err();
                let dyn_err: &(dyn std::error::Error + 'static) = &err;
                compare_errors(Some(dyn_err), &[
                    ErrorKind::PathReadMeta.into(),
                    ErrorKind::MountVersion.into(),
                ]);
                continue;
            }
            // KV2
            let res = clients.vaku.path_read_meta(&p).await;
            let err_ref: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);
            if tt.want_nil {
                assert!(res.unwrap().is_none());
            } else if tt.want_err.is_empty() {
                let meta = res.unwrap().unwrap();
                assert!(meta.current_version >= 1);
                assert!(!meta.versions.is_empty());
            }
        }
    }
}

//! Ports `api/path_update_test.go`.

use crate::common::{seeded_prefixes, shared_clients, MOUNTLESS};
use crate::skip_if_no_docker;
use serde_json::{json, Map, Value};
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;

fn mk(kvs: &[(&str, &str)]) -> Option<Map<String, Value>> {
    let mut m = Map::new();
    for (k, v) in kvs { m.insert((*k).to_string(), json!(*v)); }
    Some(m)
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_update() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        data: Option<Map<String, Value>>,
        want_data: Option<Map<String, Value>>,
        want_err: Vec<ErrMatch>,
        no_readback: bool,
    }
    let cases = vec![
        Case {
            give: "newpath",
            data: mk(&[("0", "1")]),
            want_data: mk(&[("0", "1")]),
            want_err: vec![],
            no_readback: false,
        },
        Case {
            give: "0/1",
            data: mk(&[("100", "101")]),
            want_data: mk(&[("2", "3"), ("100", "101")]),
            want_err: vec![],
            no_readback: false,
        },
        Case {
            give: "nildata",
            data: None,
            want_data: None,
            want_err: vec![ErrorKind::PathUpdate.into(), ErrorKind::NilData.into()],
            no_readback: true,
        },
        Case {
            give: "0/4/5",
            data: None,
            want_data: mk(&[("6", "7")]),
            want_err: vec![ErrorKind::PathUpdate.into(), ErrorKind::NilData.into()],
            no_readback: false,
        },
        Case {
            give: MOUNTLESS,
            data: mk(&[("0", "1")]),
            want_data: None,
            want_err: vec![
                ErrorKind::PathUpdate.into(),
                ErrorKind::PathRead.into(),
                ErrorKind::RewritePath.into(),
                ErrorKind::MountInfo.into(),
                ErrorKind::NoMount.into(),
            ],
            no_readback: true,
        },
        Case {
            give: "error/write/inject",
            data: mk(&[("0", "1")]),
            want_data: None,
            want_err: vec![
                ErrorKind::PathUpdate.into(),
                ErrorKind::PathWrite.into(),
                ErrorKind::VaultWrite.into(),
            ],
            no_readback: true,
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.path_update(&p, tt.data.clone()).await;
            let err_ref: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);

            if !tt.no_readback {
                let read = clients.clean.path_read(&p).await.unwrap();
                assert_eq!(read, tt.want_data);
            }
        }
    }
}

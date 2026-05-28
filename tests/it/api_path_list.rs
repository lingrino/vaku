//! Ports `api/path_list_test.go`.

use crate::common::{seeded_prefixes, shared_clients, MOUNTLESS};
use crate::skip_if_no_docker;
use vaku::api::client::Client;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_list};
use vaku::api::logical::VaultHttpClient;
use std::sync::Arc;

struct Case {
    give: &'static str,
    want: Vec<&'static str>,
    want_err: Vec<ErrMatch>,
}

fn cases() -> Vec<Case> {
    vec![
        Case {
            give: "0",
            want: vec!["1", "4/"],
            want_err: vec![],
        },
        Case {
            give: "0/4/13/24",
            want: vec!["25/"],
            want_err: vec![],
        },
        Case {
            give: "emptypath",
            want: vec![],
            want_err: vec![],
        },
        Case {
            give: MOUNTLESS,
            want: vec![],
            want_err: vec![
                ErrorKind::PathList.into(),
                ErrorKind::RewritePath.into(),
                ErrorKind::MountInfo.into(),
                ErrorKind::NoMount.into(),
            ],
        },
        Case {
            give: "error/list/inject",
            want: vec![],
            want_err: vec![ErrorKind::PathList.into(), ErrorKind::VaultList.into()],
        },
        Case {
            give: "nildata/list/inject",
            want: vec![],
            want_err: vec![],
        },
        Case {
            give: "nilkeys/list/inject",
            want: vec![],
            want_err: vec![ErrorKind::PathList.into(), ErrorKind::DecodeSecret.into()],
        },
        Case {
            give: "intkeys/list/inject",
            want: vec![],
            want_err: vec![ErrorKind::PathList.into(), ErrorKind::DecodeSecret.into()],
        },
        Case {
            give: "listintkeys/list/inject",
            want: vec![],
            want_err: vec![ErrorKind::PathList.into(), ErrorKind::DecodeSecret.into()],
        },
    ]
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_list() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for tt in cases() {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let result = clients.vaku.path_list(&p).await;

            let mut got = match &result {
                Ok(v) => v.clone(),
                Err(_) => Vec::new(),
            };
            trim_prefix_list(&mut got, &prefix);
            let want: Vec<String> = tt.want.iter().map(|s| s.to_string()).collect();
            assert_eq!(got, want, "give={} prefix={}", tt.give, prefix);

            let err_ref: Option<&(dyn std::error::Error + 'static)> = match result.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_list_ignore_errors() {
    skip_if_no_docker!();
    let prefix = seeded_prefixes("error/list/inject").await.into_iter().next().unwrap();
    let server = &crate::common::seeds::SERVERS.src;
    let http = Arc::new(VaultHttpClient::new(&server.addr, &server.token, None).unwrap());
    let injector = Arc::new(crate::common::injector::LogicalInjector::new(http, false));
    let client = Client::builder()
        .with_logical(injector)
        .with_ignore_access_errors(true)
        .build()
        .unwrap();

    let p = path_join(&[&prefix, "error/list/inject"]);
    let got = client.path_list(&p).await.expect("should ignore error");
    assert!(got.is_empty(), "got {got:?}");
}

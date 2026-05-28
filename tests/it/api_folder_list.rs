//! Ports `api/folder_list_test.go`.

use crate::common::{seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_list};

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_list() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        want: Vec<&'static str>,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case {
            give: "0/1",
            want: vec![],
            want_err: vec![],
        },
        Case {
            give: "0/4/",
            want: vec!["5", "8", "13/14", "13/17", "13/24/25/26/27"],
            want_err: vec![],
        },
        Case {
            give: "error/list/inject",
            want: vec![],
            want_err: vec![
                ErrorKind::FolderList.into(),
                ErrorKind::FolderListChan.into(),
                ErrorKind::PathList.into(),
                ErrorKind::VaultList.into(),
            ],
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.folder_list(&p).await;
            let mut got = res.as_ref().cloned().unwrap_or_default();
            trim_prefix_list(&mut got, &prefix);
            let mut got_sorted = got.clone();
            got_sorted.sort();
            let mut want: Vec<String> = tt.want.iter().map(|s| s.to_string()).collect();
            want.sort();
            assert_eq!(got_sorted, want, "give={} prefix={}", tt.give, prefix);

            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);
        }
    }
}

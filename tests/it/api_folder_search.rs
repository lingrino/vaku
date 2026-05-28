//! Ports `api/folder_search_test.go`.

use crate::common::{seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_list};

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_search() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        search: &'static str,
        want: Vec<&'static str>,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case { give: "0", search: "notfound", want: vec![], want_err: vec![] },
        Case { give: "0/4/13/24", search: "7", want: vec![], want_err: vec![] },
        Case { give: "0/4/13", search: "3", want: vec!["17"], want_err: vec![] },
        Case {
            give: "0/4", search: "2",
            want: vec!["8", "13/17", "13/24/25/26/27"],
            want_err: vec![],
        },
        Case {
            give: "0/4/error/read/inject", search: "aaaaaaaaa", want: vec![],
            want_err: vec![
                ErrorKind::FolderSearch.into(),
                ErrorKind::FolderRead.into(),
                ErrorKind::FolderReadChan.into(),
                ErrorKind::PathRead.into(),
                ErrorKind::VaultRead.into(),
            ],
        },
        Case {
            give: "0/4/funcdata/read/inject", search: "aaaaaaaaa", want: vec![],
            want_err: vec![ErrorKind::FolderSearch.into(), ErrorKind::JsonMarshal.into()],
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.folder_search(&p, tt.search).await;
            let mut got = res.as_ref().cloned().unwrap_or_default();
            trim_prefix_list(&mut got, &prefix);
            got.sort();
            let mut want: Vec<String> = tt.want.iter().map(|s| s.to_string()).collect();
            want.sort();
            assert_eq!(got, want, "give={} prefix={}", tt.give, prefix);

            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None, Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);
        }
    }
}

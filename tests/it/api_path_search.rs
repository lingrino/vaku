//! Ports `api/path_search_test.go`.

use crate::common::{seeded_prefixes, shared_clients, MOUNTLESS};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_search() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        give: &'static str,
        search: &'static str,
        want: bool,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case { give: "0/1", search: "2", want: true, want_err: vec![] },
        Case { give: "0/4/5", search: "7", want: true, want_err: vec![] },
        Case { give: "0/4/8", search: "13", want: false, want_err: vec![] },
        Case { give: "0/4/13/17", search: "9", want: true, want_err: vec![] },
        Case { give: "fake", search: "searchstring", want: false, want_err: vec![] },
        Case { give: "fakeempty", search: "", want: false, want_err: vec![] },
        Case {
            give: MOUNTLESS,
            search: "searchstring",
            want: false,
            want_err: vec![
                ErrorKind::PathSearch.into(),
                ErrorKind::PathRead.into(),
                ErrorKind::RewritePath.into(),
                ErrorKind::MountInfo.into(),
                ErrorKind::NoMount.into(),
            ],
        },
        Case {
            give: "error/read/inject",
            search: "searchstring",
            want: false,
            want_err: vec![
                ErrorKind::PathSearch.into(),
                ErrorKind::PathRead.into(),
                ErrorKind::VaultRead.into(),
            ],
        },
        Case {
            give: "funcdata/read/inject",
            search: "searchstring",
            want: false,
            want_err: vec![ErrorKind::PathSearch.into(), ErrorKind::JsonMarshal.into()],
        },
    ];

    for tt in cases {
        for prefix in seeded_prefixes(tt.give).await {
            let p = path_join(&[&prefix, tt.give]);
            let res = clients.vaku.path_search(&p, tt.search).await;

            let got = res.as_ref().copied().unwrap_or(false);
            assert_eq!(got, tt.want, "give={} prefix={}", tt.give, prefix);

            let err_ref: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);
        }
    }
}

//! Ports `api/path_copy_test.go`.

use crate::common::{seeded_prefix_product, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_copy() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        src: &'static str,
        dst: &'static str,
        want_err: Vec<ErrMatch>,
        nil_dst: bool,
    }
    let cases = vec![
        Case {
            src: "0/1",
            dst: "copy/0/1",
            want_err: vec![],
            nil_dst: false,
        },
        Case {
            src: "0/1",
            dst: "0/4/5",
            want_err: vec![],
            nil_dst: false,
        },
        Case {
            src: "0/4/8/error/read/inject",
            dst: "copy/readerror",
            want_err: vec![
                ErrorKind::PathCopy.into(),
                ErrorKind::PathRead.into(),
                ErrorKind::VaultRead.into(),
            ],
            nil_dst: true,
        },
        Case {
            src: "0/4/8",
            dst: "copy/writeerror/error/write/inject",
            want_err: vec![
                ErrorKind::PathCopy.into(),
                ErrorKind::PathWrite.into(),
                ErrorKind::VaultWrite.into(),
            ],
            nil_dst: true,
        },
    ];

    for tt in cases {
        for (psrc, pdst) in seeded_prefix_product().await {
            let src = path_join(&[&psrc, tt.src]);
            let dst = path_join(&[&pdst, tt.dst]);
            let res = clients.vaku.path_copy(&src, &dst).await;
            let err_ref: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None,
                Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);

            let read_src = clients.clean.path_read(&src).await.unwrap();
            let read_dst = clients
                .clean
                .as_destination()
                .path_read(&dst)
                .await
                .unwrap();
            if tt.nil_dst {
                assert!(read_dst.is_none());
            } else {
                assert_eq!(read_src, read_dst);
            }
        }
    }
}

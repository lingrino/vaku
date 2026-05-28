//! Ports `api/folder_copy_test.go`.

use crate::common::{seeded_prefix_product, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_map};

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_copy() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        src: &'static str,
        dst: &'static str,
        want_err: Vec<ErrMatch>,
        nil_dst: bool,
    }
    let cases = vec![
        Case { src: "0/1", dst: "copy/0/1", want_err: vec![], nil_dst: false },
        Case { src: "0", dst: "copy/0", want_err: vec![], nil_dst: false },
        Case {
            src: "0/4/13/24/25/26/error/read/inject", dst: "copy/0/4/13/24/25/26",
            want_err: vec![
                ErrorKind::FolderCopy.into(),
                ErrorKind::FolderRead.into(),
                ErrorKind::FolderReadChan.into(),
                ErrorKind::PathRead.into(),
                ErrorKind::VaultRead.into(),
            ],
            nil_dst: true,
        },
        Case {
            src: "0/4/13/24/25/26", dst: "copy/0/4/13/24/25/26/error/write/inject",
            want_err: vec![
                ErrorKind::FolderCopy.into(),
                ErrorKind::FolderWrite.into(),
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
            let res = clients.vaku.folder_copy(&src, &dst).await;
            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None, Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);

            let mut read_src = clients.clean.folder_read(&src).await.unwrap().unwrap_or_default();
            let mut read_dst = clients.clean.as_destination().folder_read(&dst).await.unwrap().unwrap_or_default();
            if tt.nil_dst {
                assert!(read_dst.is_empty());
            } else {
                trim_prefix_map(&mut read_src, &psrc);
                trim_prefix_map(&mut read_dst, &pdst);
                assert_eq!(read_src, read_dst);
            }
        }
    }
}

//! Ports `api/folder_copy_all_versions_test.go`.

use crate::common::{seeded_prefix_product, seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_map};

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_copy_all_versions() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        src: &'static str,
        dst: &'static str,
        want_err: Vec<ErrMatch>,
        nil_dst: bool,
    }
    let cases = vec![
        Case { src: "0/1", dst: "copyallversions/0/1", want_err: vec![], nil_dst: false },
        Case { src: "0", dst: "copyallversions/0", want_err: vec![], nil_dst: false },
        Case {
            src: "0/4/13/24/25/26/error/list/inject",
            dst: "copyallversions/error/list",
            want_err: vec![
                ErrorKind::FolderCopyAllVersions.into(),
                ErrorKind::FolderListChan.into(),
                ErrorKind::PathList.into(),
                ErrorKind::VaultList.into(),
            ],
            nil_dst: true,
        },
    ];

    for tt in cases {
        for (psrc, pdst) in seeded_prefix_product().await {
            if psrc.starts_with("kv1/") || pdst.starts_with("kv1/") { continue; }
            let src = path_join(&[&psrc, tt.src]);
            let dst = path_join(&[&pdst, tt.dst]);
            let res = clients.vaku.folder_copy_all_versions(&src, &dst).await;
            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None, Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);

            if tt.nil_dst {
                let read_dst = clients.clean.as_destination().folder_read(&dst).await.unwrap().unwrap_or_default();
                assert!(read_dst.is_empty());
            } else if tt.want_err.is_empty() {
                let mut read_src = clients.clean.folder_read(&src).await.unwrap().unwrap_or_default();
                let mut read_dst = clients.clean.as_destination().folder_read(&dst).await.unwrap().unwrap_or_default();
                trim_prefix_map(&mut read_src, &psrc);
                trim_prefix_map(&mut read_dst, &pdst);
                assert_eq!(read_src, read_dst);
            }
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_copy_all_versions_kv1() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv1/") {
            let err = clients.vaku
                .folder_copy_all_versions(&path_join(&[&prefix, "0/1"]), "kv2/copyallversions/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::FolderCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv2/") {
            let err = clients.vaku
                .folder_copy_all_versions(&path_join(&[&prefix, "0/1"]), "kv1/copyallversions/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::FolderCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
}

//! Ports `api/folder_move_all_versions_test.go`.

use crate::common::{seeded_prefix_product, seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::{path_join, trim_prefix_map};

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_move_all_versions() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        src: &'static str,
        dst: &'static str,
        want_err: Vec<ErrMatch>,
        nil_src: bool,
        nil_dst: bool,
    }
    let cases = vec![
        Case { src: "0/1", dst: "moveallversions/0/1", want_err: vec![], nil_src: true, nil_dst: false },
        Case { src: "0", dst: "moveallversions/0", want_err: vec![], nil_src: true, nil_dst: false },
        Case {
            src: "0/4/13/24/25/26/error/list/inject", dst: "moveallversions/error/list",
            want_err: vec![
                ErrorKind::FolderMoveAllVersions.into(),
                ErrorKind::FolderCopyAllVersions.into(),
                ErrorKind::FolderListChan.into(),
                ErrorKind::PathList.into(),
                ErrorKind::VaultList.into(),
            ],
            nil_src: false, nil_dst: true,
        },
        Case {
            src: "0/4/13/24/25/26/error/delete/inject", dst: "moveallversions/error/delete",
            want_err: vec![
                ErrorKind::FolderMoveAllVersions.into(),
                ErrorKind::FolderDeleteMeta.into(),
                ErrorKind::PathDeleteMeta.into(),
                ErrorKind::VaultDelete.into(),
            ],
            nil_src: false, nil_dst: false,
        },
    ];

    for tt in cases {
        for (psrc, pdst) in seeded_prefix_product().await {
            if psrc.starts_with("kv1/") || pdst.starts_with("kv1/") { continue; }
            let src = path_join(&[&psrc, tt.src]);
            let dst = path_join(&[&pdst, tt.dst]);
            let mut orig = clients.clean.folder_read(&src).await.unwrap().unwrap_or_default();
            trim_prefix_map(&mut orig, &src);

            let res = clients.vaku.folder_move_all_versions(&src, &dst).await;
            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None, Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);

            let mut read_src = clients.clean.folder_read(&src).await.unwrap().unwrap_or_default();
            let mut read_dst = clients.clean.as_destination().folder_read(&dst).await.unwrap().unwrap_or_default();
            trim_prefix_map(&mut read_src, &src);
            trim_prefix_map(&mut read_dst, &dst);
            if tt.nil_src { assert!(read_src.is_empty()); }
            else { assert_eq!(read_src, orig); }
            if tt.nil_dst { assert!(read_dst.is_empty()); }
            else { assert_eq!(read_dst, orig); }
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_folder_move_all_versions_kv1() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv1/") {
            let err = clients.vaku
                .folder_move_all_versions(&path_join(&[&prefix, "0/1"]), "kv2/moveallversions/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::FolderMoveAllVersions.into(),
                ErrorKind::FolderCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv2/") {
            let err = clients.vaku
                .folder_move_all_versions(&path_join(&[&prefix, "0/1"]), "kv1/moveallversions/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::FolderMoveAllVersions.into(),
                ErrorKind::FolderCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
}

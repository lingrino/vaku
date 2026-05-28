//! Ports `api/path_move_all_versions_test.go`.

use crate::common::{seeded_prefix_product, seeded_prefixes, shared_clients};
use crate::skip_if_no_docker;
use serde_json::{json, Map};
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::helpers::path_join;

fn mk(kvs: &[(&str, &str)]) -> Option<Map<String, serde_json::Value>> {
    let mut m = Map::new();
    for (k, v) in kvs { m.insert((*k).to_string(), json!(*v)); }
    Some(m)
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_move_all_versions() {
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
        Case { src: "0/1", dst: "move/allversions/0/1", want_err: vec![], nil_src: true, nil_dst: false },
        Case { src: "0/4/5", dst: "move/allversions/different", want_err: vec![], nil_src: true, nil_dst: false },
        Case { src: "fake/nonexistent", dst: "move/allversions/fake", want_err: vec![], nil_src: true, nil_dst: true },
        Case {
            src: "error/read/inject",
            dst: "move/allversions/readerror",
            want_err: vec![
                ErrorKind::PathMoveAllVersions.into(),
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::PathReadMeta.into(),
                ErrorKind::VaultRead.into(),
            ],
            nil_src: true, nil_dst: true,
        },
        Case {
            src: "0/4/8",
            dst: "move/allversions/error/write/inject",
            want_err: vec![
                ErrorKind::PathMoveAllVersions.into(),
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::PathWrite.into(),
                ErrorKind::VaultWrite.into(),
            ],
            nil_src: false, nil_dst: true,
        },
        Case {
            src: "0/4/13/14/error/delete/inject",
            dst: "move/allversions/deleteerror",
            want_err: vec![
                ErrorKind::PathMoveAllVersions.into(),
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
            let orig = clients.clean.path_read(&src).await.unwrap();
            let res = clients.vaku.path_move_all_versions(&src, &dst).await;
            let er: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None, Err(e) => Some(e),
            };
            compare_errors(er, &tt.want_err);

            let read_src = clients.clean.path_read(&src).await.unwrap();
            let read_dst = clients.clean.as_destination().path_read(&dst).await.unwrap();
            if tt.nil_src { assert!(read_src.is_none()); }
            else { assert_eq!(read_src, orig); }
            if tt.nil_dst { assert!(read_dst.is_none()); }
            else { assert_eq!(read_dst, orig); }
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_move_all_versions_kv1() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv1/") {
            let err = clients.vaku
                .path_move_all_versions(&path_join(&[&prefix, "0/1"]), "kv2/move/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::PathMoveAllVersions.into(),
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv2/") {
            let err = clients.vaku
                .path_move_all_versions(&path_join(&[&prefix, "0/1"]), "kv1/move/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::PathMoveAllVersions.into(),
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_move_all_versions_with_multiple() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("movemulti").await {
        if !prefix.starts_with("kv2/") { continue; }
        let src = path_join(&[&prefix, "movemulti/test"]);
        let dst = path_join(&[&prefix, "move/movemulti/test"]);
        let v1 = mk(&[("version", "1"), ("data", "first")]).unwrap();
        let v2 = mk(&[("version", "2"), ("data", "second")]).unwrap();
        let v3 = mk(&[("version", "3"), ("data", "third")]).unwrap();
        clients.clean.path_write(&src, Some(v1.clone())).await.unwrap();
        clients.clean.path_write(&src, Some(v2.clone())).await.unwrap();
        clients.clean.path_write(&src, Some(v3.clone())).await.unwrap();

        clients.vaku.path_move_all_versions(&src, &dst).await.unwrap();

        assert!(clients.clean.path_read(&src).await.unwrap().is_none());
        let dst_meta = clients.clean.as_destination().path_read_meta(&dst).await.unwrap().unwrap();
        assert_eq!(dst_meta.versions.len(), 3);
    }
}

//! Ports `api/path_copy_all_versions_test.go`.

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
async fn test_path_copy_all_versions() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    struct Case {
        src: &'static str,
        dst: &'static str,
        want_err: Vec<ErrMatch>,
        nil_dst: bool,
    }
    let cases = vec![
        Case { src: "0/1", dst: "copy/allversions/0/1", want_err: vec![], nil_dst: false },
        Case { src: "0/4/5", dst: "copy/allversions/different", want_err: vec![], nil_dst: false },
        Case { src: "fake/nonexistent", dst: "copy/allversions/fake", want_err: vec![], nil_dst: true },
        Case {
            src: "0/4/8/error/read/inject",
            dst: "copy/allversions/readerror",
            want_err: vec![
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::PathReadMeta.into(),
                ErrorKind::VaultRead.into(),
            ],
            nil_dst: true,
        },
        Case {
            src: "0/4/8",
            dst: "copy/allversions/error/write/inject",
            want_err: vec![
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::PathWrite.into(),
                ErrorKind::VaultWrite.into(),
            ],
            nil_dst: true,
        },
    ];

    for tt in cases {
        for (psrc, pdst) in seeded_prefix_product().await {
            if psrc.starts_with("kv1/") || pdst.starts_with("kv1/") { continue; }
            let src = path_join(&[&psrc, tt.src]);
            let dst = path_join(&[&pdst, tt.dst]);
            let res = clients.vaku.path_copy_all_versions(&src, &dst).await;
            let err_ref: Option<&(dyn std::error::Error + 'static)> = match res.as_ref() {
                Ok(_) => None, Err(e) => Some(e),
            };
            compare_errors(err_ref, &tt.want_err);

            if tt.nil_dst {
                let read_dst = clients.clean.as_destination().path_read(&dst).await.unwrap();
                assert!(read_dst.is_none());
            } else if tt.want_err.is_empty() {
                let read_src = clients.clean.path_read(&src).await.unwrap();
                let read_dst = clients.clean.as_destination().path_read(&dst).await.unwrap();
                assert_eq!(read_src, read_dst);

                let src_meta = clients.clean.path_read_meta(&src).await.unwrap();
                let dst_meta = clients.clean.as_destination().path_read_meta(&dst).await.unwrap();
                if let (Some(s), Some(d)) = (src_meta, dst_meta) {
                    assert_eq!(s.versions.len(), d.versions.len());
                }
            }
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_copy_all_versions_kv1() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv1/") {
            let err = clients.vaku
                .path_copy_all_versions(&path_join(&[&prefix, "0/1"]), "kv2/copy/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
    for prefix in seeded_prefixes("0/1").await {
        if prefix.starts_with("kv2/") {
            let err = clients.vaku
                .path_copy_all_versions(&path_join(&[&prefix, "0/1"]), "kv1/copy/dst")
                .await.unwrap_err();
            let de: &(dyn std::error::Error + 'static) = &err;
            compare_errors(Some(de), &[
                ErrorKind::PathCopyAllVersions.into(),
                ErrorKind::MountVersion.into(),
            ]);
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_copy_all_versions_with_deleted_versions() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("multiversion").await {
        if !prefix.starts_with("kv2/") { continue; }
        let src = path_join(&[&prefix, "multiversion/test"]);
        let dst = path_join(&[&prefix, "copy/multiversion/test"]);

        let v1 = mk(&[("version", "1"), ("data", "first")]).unwrap();
        let v2 = mk(&[("version", "2"), ("data", "second")]).unwrap();
        let v3 = mk(&[("version", "3"), ("data", "third")]).unwrap();
        clients.clean.path_write(&src, Some(v1.clone())).await.unwrap();
        clients.clean.path_write(&src, Some(v2.clone())).await.unwrap();
        clients.clean.path_write(&src, Some(v3.clone())).await.unwrap();

        let meta = clients.clean.path_read_meta(&src).await.unwrap().unwrap();
        assert_eq!(meta.versions.len(), 3);

        clients.vaku.path_copy_all_versions(&src, &dst).await.unwrap();

        let dst_meta = clients.clean.as_destination().path_read_meta(&dst).await.unwrap().unwrap();
        assert_eq!(dst_meta.versions.len(), 3);

        let d1 = clients.clean.as_destination().path_read_version(&dst, 1).await.unwrap().unwrap();
        assert_eq!(d1, v1);
        let d2 = clients.clean.as_destination().path_read_version(&dst, 2).await.unwrap().unwrap();
        assert_eq!(d2, v2);
        let d3 = clients.clean.as_destination().path_read_version(&dst, 3).await.unwrap().unwrap();
        assert_eq!(d3, v3);
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_path_copy_all_versions_with_destroyed() {
    skip_if_no_docker!();
    let clients = shared_clients().await;
    for prefix in seeded_prefixes("destroyed").await {
        if !prefix.starts_with("kv2/") { continue; }
        let src = path_join(&[&prefix, "destroyed/test"]);
        let dst = path_join(&[&prefix, "copy/destroyed/test"]);

        let v1 = mk(&[("version", "1")]).unwrap();
        let v2 = mk(&[("version", "2")]).unwrap();
        let v3 = mk(&[("version", "3")]).unwrap();
        clients.clean.path_write(&src, Some(v1.clone())).await.unwrap();
        clients.clean.path_write(&src, Some(v2.clone())).await.unwrap();
        clients.clean.path_write(&src, Some(v3.clone())).await.unwrap();

        clients.clean.path_destroy(&src, &[2]).await.unwrap();
        let meta = clients.clean.path_read_meta(&src).await.unwrap().unwrap();
        assert!(meta.versions[&2].destroyed);

        clients.vaku.path_copy_all_versions(&src, &dst).await.unwrap();

        let dst_meta = clients.clean.as_destination().path_read_meta(&dst).await.unwrap().unwrap();
        assert_eq!(dst_meta.versions.len(), 3);

        let d1 = clients.clean.as_destination().path_read_version(&dst, 1).await.unwrap().unwrap();
        assert_eq!(d1, v1);
        let d2 = clients.clean.as_destination().path_read_version(&dst, 2).await.unwrap();
        assert!(d2.map_or(true, |m| m.is_empty()));
        let d3 = clients.clean.as_destination().path_read_version(&dst, 3).await.unwrap().unwrap();
        assert_eq!(d3, v3);
    }
}

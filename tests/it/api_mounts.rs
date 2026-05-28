//! Ports `api/mounts_test.go`.

use crate::common::{shared_clients, SERVERS};
use crate::skip_if_no_docker;
use std::sync::Arc;
use vaku::api::client::Client;
use vaku::api::error::{compare_errors, ErrMatch, ErrorKind};
use vaku::api::logical::VaultHttpClient;
use vaku::api::mount_provider::StaticMountProvider;
use vaku::api::mounts::{mount_info, rewrite_path, MountVersion, VaultOp};

#[test]
fn test_static_mount_provider() {
    let cases = [("secret/", "2"), ("kv1/", "1"), ("my/secret/", "2")];
    for (path, version) in cases {
        let provider = StaticMountProvider::new(path, version);
        let mounts = tokio::runtime::Runtime::new()
            .unwrap()
            .block_on(
                <StaticMountProvider as vaku::api::mount_provider::MountProvider>::list_mounts(
                    &provider,
                ),
            )
            .unwrap();
        assert_eq!(mounts.len(), 1);
        assert_eq!(mounts[0].path, path);
        assert_eq!(mounts[0].version, version);
        assert_eq!(mounts[0].r#type, "kv");
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_mount_info_with_static_provider() {
    struct Case {
        mount_path: &'static str,
        mount_ver: &'static str,
        query: &'static str,
        want_path: &'static str,
        want_ver: MountVersion,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case {
            mount_path: "secret/",
            mount_ver: "2",
            query: "secret/foo/bar",
            want_path: "secret/",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
        Case {
            mount_path: "kv1/",
            mount_ver: "1",
            query: "kv1/foo/bar",
            want_path: "kv1/",
            want_ver: MountVersion::Mv1,
            want_err: vec![],
        },
        Case {
            mount_path: "secret/",
            mount_ver: "2",
            query: "other/foo/bar",
            want_path: "",
            want_ver: MountVersion::Mv0,
            want_err: vec![ErrorKind::MountInfo.into(), ErrorKind::NoMount.into()],
        },
    ];

    for tt in cases {
        let provider = Arc::new(StaticMountProvider::new(tt.mount_path, tt.mount_ver));
        let res = mount_info(provider.as_ref(), tt.query).await;
        match res {
            Ok((path, ver)) => {
                assert_eq!(path, tt.want_path);
                assert_eq!(ver, tt.want_ver);
                assert!(tt.want_err.is_empty());
            }
            Err(e) => {
                let de: &(dyn std::error::Error + 'static) = &e;
                compare_errors(Some(de), &tt.want_err);
            }
        }
    }
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_mount_info() {
    skip_if_no_docker!();
    // Empty-mounts client — bogus token + addr so list_mounts errors.
    let bad = VaultHttpClient::new("http://127.0.0.1:1/", "bad", None).unwrap();
    let client = Client::builder().with_vault_client(bad).build().unwrap();
    let res = mount_info(client.src().mount_provider.as_ref(), "kv0").await;
    let err = res.unwrap_err();
    let de: &(dyn std::error::Error + 'static) = &err;
    compare_errors(
        Some(de),
        &[ErrorKind::MountInfo.into(), ErrorKind::ListMounts.into()],
    );

    // Live container — use the clean shared client (its mount provider hits
    // the real sys/mounts).
    let _ = shared_clients().await;
    let server = &SERVERS.src;
    let http = VaultHttpClient::new(&server.addr, &server.token, None).unwrap();
    let cl = Client::builder().with_vault_client(http).build().unwrap();

    struct Case {
        give: &'static str,
        want_path: &'static str,
        want_ver: MountVersion,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case {
            give: "nomount",
            want_path: "",
            want_ver: MountVersion::Mv0,
            want_err: vec![ErrorKind::MountInfo.into(), ErrorKind::NoMount.into()],
        },
        Case {
            give: "sys/",
            want_path: "sys/",
            want_ver: MountVersion::Mv0,
            want_err: vec![],
        },
        Case {
            give: "kv1/",
            want_path: "kv1/",
            want_ver: MountVersion::Mv1,
            want_err: vec![],
        },
        Case {
            give: "kv2/",
            want_path: "kv2/",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
    ];

    for tt in cases {
        let res = mount_info(cl.src().mount_provider.as_ref(), tt.give).await;
        match res {
            Ok((path, ver)) => {
                assert_eq!(path, tt.want_path, "give={}", tt.give);
                assert_eq!(ver, tt.want_ver);
                assert!(tt.want_err.is_empty());
            }
            Err(e) => {
                let de: &(dyn std::error::Error + 'static) = &e;
                compare_errors(Some(de), &tt.want_err);
            }
        }
    }
}

#[test]
fn test_mount_string_to_version() {
    use MountVersion::*;
    assert_eq!(MountVersion::parse("---"), Mv0);
    assert_eq!(MountVersion::parse("0"), Mv0);
    assert_eq!(MountVersion::parse("1"), Mv1);
    assert_eq!(MountVersion::parse("2"), Mv2);
    assert_eq!(MountVersion::parse("3"), Other(3));
}

#[tokio::test(flavor = "multi_thread", worker_threads = 4)]
async fn test_rewrite_path() {
    skip_if_no_docker!();
    let _ = shared_clients().await;
    let server = &SERVERS.src;
    let http = VaultHttpClient::new(&server.addr, &server.token, None).unwrap();
    let cl = Client::builder().with_vault_client(http).build().unwrap();

    struct Case {
        give: &'static str,
        op: VaultOp,
        want_path: &'static str,
        want_ver: MountVersion,
        want_err: Vec<ErrMatch>,
    }
    let cases = vec![
        Case {
            give: "nomount",
            op: VaultOp::Read,
            want_path: "",
            want_ver: MountVersion::Mv0,
            want_err: vec![
                ErrorKind::RewritePath.into(),
                ErrorKind::MountInfo.into(),
                ErrorKind::NoMount.into(),
            ],
        },
        Case {
            give: "kv1/a/b/c",
            op: VaultOp::List,
            want_path: "kv1/a/b/c",
            want_ver: MountVersion::Mv1,
            want_err: vec![],
        },
        Case {
            give: "kv2/a/b/c",
            op: VaultOp::List,
            want_path: "kv2/metadata/a/b/c",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
        Case {
            give: "kv1/a/b/c",
            op: VaultOp::Read,
            want_path: "kv1/a/b/c",
            want_ver: MountVersion::Mv1,
            want_err: vec![],
        },
        Case {
            give: "kv2/a/b/c",
            op: VaultOp::Read,
            want_path: "kv2/data/a/b/c",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
        Case {
            give: "kv1/a/b/c",
            op: VaultOp::Write,
            want_path: "kv1/a/b/c",
            want_ver: MountVersion::Mv1,
            want_err: vec![],
        },
        Case {
            give: "kv2/a/b/c",
            op: VaultOp::Write,
            want_path: "kv2/data/a/b/c",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
        Case {
            give: "kv1/a/b/c",
            op: VaultOp::Delete,
            want_path: "kv1/a/b/c",
            want_ver: MountVersion::Mv1,
            want_err: vec![],
        },
        Case {
            give: "kv2/a/b/c",
            op: VaultOp::Delete,
            want_path: "kv2/data/a/b/c",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
        Case {
            give: "kv1/a/b/c",
            op: VaultOp::Destroy,
            want_path: "",
            want_ver: MountVersion::Mv1,
            want_err: vec![ErrorKind::MountVersion.into()],
        },
        Case {
            give: "kv2/a/b/c",
            op: VaultOp::Destroy,
            want_path: "kv2/destroy/a/b/c",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
        Case {
            give: "kv1/a/b/c",
            op: VaultOp::DeleteMeta,
            want_path: "",
            want_ver: MountVersion::Mv1,
            want_err: vec![ErrorKind::MountVersion.into()],
        },
        Case {
            give: "kv2/a/b/c",
            op: VaultOp::DeleteMeta,
            want_path: "kv2/metadata/a/b/c",
            want_ver: MountVersion::Mv2,
            want_err: vec![],
        },
    ];

    for tt in cases {
        let res = rewrite_path(cl.src().mount_provider.as_ref(), tt.give, tt.op).await;
        match res {
            Ok((path, ver)) => {
                assert_eq!(path, tt.want_path, "give={}", tt.give);
                assert_eq!(ver, tt.want_ver);
                assert!(tt.want_err.is_empty());
            }
            Err(e) => {
                let de: &(dyn std::error::Error + 'static) = &e;
                compare_errors(Some(de), &tt.want_err);
            }
        }
    }
}

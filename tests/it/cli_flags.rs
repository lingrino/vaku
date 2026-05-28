//! Tests for flag validation.

use crate::cli_common::run_cli_with_mock;

fn run(args: &[&str]) -> (i32, String, String) {
    run_cli_with_mock(args)
}

#[test]
fn invalid_format_rejected() {
    let (code, _, err) = run(&["--format=invalid", "path", "list", "foo"]);
    assert_eq!(code, 1);
    assert!(err.contains("format must be one of"), "err={err}");
}

#[test]
fn invalid_workers_rejected() {
    let (code, _, err) = run(&["path", "list", "--workers=0", "foo"]);
    assert_eq!(code, 1);
    assert!(err.contains("workers must be >= 1"), "err={err}");
}

#[test]
fn mount_version_requires_path() {
    let (code, _, err) = run(&["path", "list", "--mount-version=1", "foo"]);
    assert_eq!(code, 1);
    assert!(
        err.contains("mount-version requires --mount-path"),
        "err={err}"
    );
}

#[test]
fn invalid_mount_version() {
    let (code, _, err) = run(&[
        "path",
        "list",
        "--mount-path=secret/",
        "--mount-version=3",
        "foo",
    ]);
    assert_eq!(code, 1);
    assert!(
        err.contains("mount-version must be one of: 1|2"),
        "err={err}"
    );
}

#[test]
fn source_mount_version_requires_path() {
    let (code, _, err) = run(&["path", "list", "--mount-version-source=1", "foo"]);
    assert_eq!(code, 1);
    assert!(
        err.contains("mount-version-source requires --mount-path-source"),
        "err={err}"
    );
}

#[test]
fn dst_mount_version_requires_path() {
    let (code, _, err) = run(&["path", "list", "--mount-version-destination=1", "foo"]);
    assert_eq!(code, 1);
    assert!(
        err.contains("mount-version-destination requires --mount-path-destination"),
        "err={err}"
    );
}

#[test]
fn valid_mount_flags() {
    let (code, _, _) = run(&[
        "path",
        "list",
        "--mount-path=secret/",
        "--mount-version=2",
        "foo",
    ]);
    assert_eq!(code, 0);
}

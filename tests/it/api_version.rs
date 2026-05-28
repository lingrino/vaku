//! Ports `api/version_test.go`.

#[test]
fn test_version() {
    assert_eq!(vaku::api::version::version(), "3.0.0");
}

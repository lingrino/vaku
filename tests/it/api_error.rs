//! Ports `api/error_test.go`.

use std::error::Error as StdError;
use std::fmt;
use vaku::api::error::{compare_errors, ErrMatch, Error, ErrorKind};

#[derive(Debug)]
struct StaticErr(&'static str);
impl fmt::Display for StaticErr {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.0)
    }
}
impl StdError for StaticErr {}

fn inject() -> Box<dyn StdError + Send + Sync> {
    Box::new(StaticErr("injected error"))
}

#[test]
fn new_wrap_err_nil_all() {
    let e = Error::new(None, None, None);
    assert_eq!(e.kind(), &ErrorKind::UnknownError);
    assert_eq!(e.to_string(), "unknown error");
    assert!(e.source().is_none());
}

#[test]
fn new_wrap_err_nil_msg_and_is() {
    let e = Error::new(None, None, Some(inject()));
    assert_eq!(e.kind(), &ErrorKind::UnknownError);
    assert_eq!(e.to_string(), "unknown error: injected error");
}

#[test]
fn new_wrap_err_nil_is_and_wrap_with_msg() {
    let e = Error::new(Some("random error".to_string()), None, None);
    match e.kind() {
        ErrorKind::Custom(s) => assert_eq!(s, "random error"),
        other => panic!("expected Custom, got {other:?}"),
    }
    assert_eq!(e.to_string(), "random error");
}

#[test]
fn new_wrap_err_nil_msg_and_wrap() {
    let e = Error::new(None, Some(ErrorKind::Custom("injected error".into())), None);
    match e.kind() {
        ErrorKind::Custom(s) => assert_eq!(s, "injected error"),
        other => panic!("expected Custom, got {other:?}"),
    }
    assert_eq!(e.to_string(), "injected error");
}

#[test]
fn new_wrap_err_msg_and_nil_wrap() {
    let e = Error::new(
        Some("random error".to_string()),
        Some(ErrorKind::Custom("injected error".into())),
        None,
    );
    assert_eq!(e.to_string(), "random error: injected error");
}

#[test]
fn new_wrap_err_standard_error() {
    let e = Error::new(
        Some("context here".to_string()),
        Some(ErrorKind::Custom("standard error".into())),
        Some(inject()),
    );
    assert_eq!(
        e.to_string(),
        "context here: standard error: injected error"
    );
}

#[test]
fn ctx_err_nil() {
    assert!(Error::ctx(None).is_none());
}

#[test]
fn ctx_err_wraps() {
    let e = Error::ctx(Some(inject())).expect("some");
    let dyn_err: &(dyn StdError + 'static) = &e;
    compare_errors(
        Some(dyn_err),
        &[
            ErrMatch::Plain(ErrorKind::Context),
            ErrMatch::Custom("injected error".to_string()),
        ],
    );
}

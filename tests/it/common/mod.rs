//! Shared test harness.
//!
//! Mirrors Go's `api/main_test.go`: one or two Vault dev-mode docker
//! containers, two shared clients (one with a path-marker-based logical
//! injector for error/edge-case injection, one without), and a seeded-prefix
//! scheme that gives every test isolated mount-prefixed paths.
//!
//! All integration tests share these globals via [`shared_clients`]; the
//! Vault containers live for the duration of the test binary.

#![allow(dead_code)]

pub mod docker;
pub mod injector;
pub mod seeds;

pub use docker::*;
pub use injector::*;
pub use seeds::*;

/// Check whether Docker is usable. Tests that need a live Vault container
/// call this and `return` early if it returns `false`, so the test binary
/// stays useful on machines without Docker.
pub fn docker_available() -> bool {
    static AVAILABLE: once_cell::sync::OnceCell<bool> = once_cell::sync::OnceCell::new();
    *AVAILABLE.get_or_init(|| {
        if std::env::var("VAKU_SKIP_LIVE_TESTS").is_ok() {
            return false;
        }
        std::process::Command::new("docker")
            .args(["info"])
            .stdout(std::process::Stdio::null())
            .stderr(std::process::Stdio::null())
            .status()
            .map(|s| s.success())
            .unwrap_or(false)
    })
}

/// Skip a test when docker isn't available.
#[macro_export]
macro_rules! skip_if_no_docker {
    () => {
        if !$crate::common::docker_available() {
            eprintln!("skipping: docker not available");
            return;
        }
    };
}

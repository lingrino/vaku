//! Library version.

/// Returns the current Vaku library version.
///
/// Used by the `vaku version` CLI subcommand.
pub fn version() -> &'static str {
    env!("CARGO_PKG_VERSION")
}

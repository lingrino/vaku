//! Vaku CLI. Mirrors the Go cobra-based CLI 1:1 in flags, output format,
//! and subcommand surface.
//!
//! Public entry point: [`execute`]. It builds the clap [`Cli`] struct,
//! dispatches, and writes output / errors through the supplied writers so
//! tests can capture stdout/stderr.

pub mod args;
pub mod client_iface;
pub mod docs;
pub mod errors;
pub mod helpers;
pub mod runner;

use std::io::Write;

/// Library-style entrypoint mirroring Go's `cmd.Execute`. Returns the exit
/// code.
pub fn execute(version: &str, args: &[String], out: &mut dyn Write, err: &mut dyn Write) -> u8 {
    runner::run(version, args, out, err)
}

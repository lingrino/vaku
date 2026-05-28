//! CLI surface for `vaku`. Phase-7 work — currently a minimal stub so the
//! binary compiles. Replaced by the full implementation later.

use std::io::Write;

pub fn execute(version: &str, args: &[String], out: &mut dyn Write, err: &mut dyn Write) -> u8 {
    if args.iter().any(|a| a == "version") {
        let _ = writeln!(out, "vaku {version}");
        return 0;
    }
    let _ = writeln!(err, "vaku CLI not yet wired up");
    1
}

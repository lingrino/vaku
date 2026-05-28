# Vaku

[![Vaku](www/assets/images/logo-vaku-sm.png?raw=true)](www/assets/logo-vaku-sm.png "Vaku")

Vaku is a CLI and Rust library for running path- and folder-based
operations on Vault's K/V secrets engine. Vaku extends the existing Vault
CLI/API by letting you run the same list / read / write / delete actions on
**folders** as well as paths, and adds copy, move, search, and bulk
delete/destroy across either source or source/destination Vault clusters.

> Vaku was originally written in Go. **Version 3 is a full rewrite in
> Rust**. The CLI surface (flags, subcommands, output) is identical; the
> public library API now lives at [docs.rs/vaku](https://docs.rs/vaku).

## Installation

### Binary

Download the latest binary for your OS/arch from the
[releases page](https://github.com/lingrino/vaku/releases) — `.tar.gz`
archives are published for Linux (x86_64, aarch64), macOS (x86_64,
aarch64), and Windows (x86_64).

### Cargo

```shell
cargo install vaku
```

### Docker

```shell
docker run ghcr.io/lingrino/vaku --help
```

### Homebrew

The Rust rewrite is not yet on the Homebrew tap. Use the binary or
`cargo install` until v3.x is published there.

## Usage

```shell
vaku --help
vaku path list secret/foo
vaku folder copy secret/old/ secret/new/
```

Full per-command docs live in [docs/cli](docs/cli/vaku.md) and are
regenerated from the binary with `vaku docs docs/cli`.

## Library

```toml
# Cargo.toml
[dependencies]
vaku = "3"
```

```rust
use std::sync::Arc;
use vaku::{Client, VaultHttpClient};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let http = VaultHttpClient::new("http://127.0.0.1:8200", "dev-root-token", None)?;
    let client = Client::builder().with_vault_client(http).build()?;
    let secrets = client.folder_read("secret/foo/").await?;
    println!("{:#?}", secrets);
    Ok(())
}
```

## Tests

Vaku is well tested. Most tests run against a live `hashicorp/vault`
Docker container.

```shell
# Full suite (requires Docker)
cargo test --all-features

# Skip the live-Vault tests
VAKU_SKIP_LIVE_TESTS=1 cargo test
```

CI runs `cargo fmt`, `cargo clippy -- -D warnings`, and `cargo test` on
every push & PR.

## Contributing

Bug reports, ideas, and PRs are all welcome. Open an issue first for
substantial changes so we can sketch the approach together.

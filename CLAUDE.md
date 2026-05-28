# Claude Code Guidelines for Vaku (Rust)

## Project Overview

Vaku is a Rust CLI and library for managing HashiCorp Vault secrets. It
provides path- and folder-level operations on the K/V secrets engine
(KV v1 and v2) including read, write, copy, move, delete, and search.

This project was originally written in Go and was rewritten in Rust for
the v3.0.0 release.

## Project Structure

```
src/
  lib.rs                # Public library re-exports
  main.rs               # `vaku` binary entry point
  api/                  # Library implementation (one file per operation)
    client.rs           # Client + ClientBuilder
    error.rs            # Error + ErrorKind + compare_errors test helper
    helpers.rs          # path_join, ensure_folder, prefix helpers
    logical.rs          # Logical trait + reqwest HTTP implementation
    mount_provider.rs   # MountProvider trait + DefaultMountProvider + StaticMountProvider
    mounts.rs           # mount_info, rewrite_path, KV1/KV2 op routing
    secret.rs           # Secret, SecretMeta types + KV2 decoders
    path_*.rs           # Single-path operations
    folder_*.rs         # Recursive folder operations
    version.rs          # Version constant
  cli/                  # CLI implementation built on clap
    args.rs             # Arg/Subcommand structs
    runner.rs           # Dispatch / token resolution / flag validation
    helpers.rs          # text/JSON renderers (Go-byte-compatible)
    docs.rs             # Markdown docs emitter
tests/it/
  main.rs               # Single integration-test binary
  common/               # Shared test harness (Vault container + injector)
  api_*.rs              # Ports of api/*_test.go
  cli_*.rs              # CLI tests using a mock ClientInterface
docs/cli/               # Auto-generated CLI docs
```

## Development Commands

```bash
# Build
cargo build              # debug
cargo build --release    # production

# Tests
cargo test               # full suite, requires Docker for live Vault
VAKU_SKIP_LIVE_TESTS=1 cargo test   # skip live tests

# Lint / format
cargo fmt --all -- --check
cargo clippy --all-targets --all-features -- -D warnings

# Regenerate docs (REQUIRED after adding/changing CLI flags)
cargo run -- docs docs/cli
```

## CI Requirements

Before pushing, ensure:
1. `cargo fmt --all -- --check` passes
2. `cargo clippy --all-targets --all-features -- -D warnings` passes
3. `cargo test` passes (or at least `VAKU_SKIP_LIVE_TESTS=1 cargo test`)
4. `cargo run -- docs docs/cli` has been run if CLI flags changed

## Code Patterns

### Adding a New API Method

1. Create `src/api/<operation>.rs`:

   ```rust
   use crate::api::client::Client;
   use crate::api::error::{Error, ErrorKind};

   impl Client {
       /// Brief doc.
       pub async fn operation_name(&self, ...) -> Result<..., Error> {
           // implementation
       }
   }
   ```

2. Add the module to `src/api/mod.rs`.
3. Add the method to `ClientInterface` in `src/cli/client_iface.rs`.
4. Create `tests/it/api_<operation>.rs` mirroring the matching Go test
   file and add it as a `mod` in `tests/it/main.rs`.
5. Add a mock impl in `tests/it/cli_common.rs::MockClient` if the CLI
   surfaces the new method.

### Adding a CLI Command/Flag

1. Add the variant / field to `src/cli/args.rs` (subcommand or flag).
   For flags that should propagate to subcommands, mark them
   `#[arg(global = true)]`.
2. Dispatch in `src/cli/runner.rs::dispatch_path` / `dispatch_folder`.
3. Add CLI tests in `tests/it/cli_basic.rs` or a new file.
4. Regenerate docs: `cargo run -- docs docs/cli`.

### Error Handling

Use `Error::wrap(msg, ErrorKind::Xxx, source)`:

```rust
return Err(Error::wrap(p, ErrorKind::PathRead, Some(Box::new(inner))));
```

Tests assert exact chains via `compare_errors`:

```rust
compare_errors(Some(err_as_dyn), &[
    ErrorKind::PathRead.into(),
    ErrorKind::VaultRead.into(),
]);
```

### Folder Operations with Workers

Use the path queue + worker pool pattern in
`src/api/folder_list.rs` as a template:
`async_channel::unbounded` for the path queue, an `AtomicUsize`
counter for "outstanding paths" (increment **before** sending),
`tokio::task::JoinSet` for workers, and a
`tokio_util::sync::CancellationToken` to short-circuit on first error.

### KV v2-Only Operations

Validate upfront via `validate_kv2`:

```rust
validate_kv2(self.src().mount_provider.as_ref(), src).await
    .map_err(|e| Error::wrap(src, ErrorKind::PathCopyAllVersions, Some(Box::new(e))))?;
```

## Common Pitfalls

1. **Forgetting to regenerate docs** — `docs/cli/*.md` must be
   re-emitted when CLI args change.
2. **Mock-driven CLI tests use `--indent-char=""`** — the test harness
   does this automatically (matches Go's `cli.flagIndent = ""`).
3. **Live tests need Docker** — `VAKU_SKIP_LIVE_TESTS=1` opts out so
   pure-function tests still run locally without Docker.
4. **Concurrency: increment-before-send** — when the folder-list
   recursion engine discovers children, it must bump the pending
   counter before pushing them onto the queue, otherwise the queue
   may close prematurely.

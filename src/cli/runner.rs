//! Dispatch & execute clap-parsed commands.

use crate::api::client::Client;
use crate::api::logical::VaultHttpClient;
use crate::api::mount_provider::StaticMountProvider;
use crate::cli::args::{FolderSub, PathFolderFlags, PathSub, TopCmd, VakuArgs};
use crate::cli::client_iface::ClientInterface;
use crate::cli::docs;
use crate::cli::errors::*;
use crate::cli::helpers::{output, Out, OutputCtx};
use clap::{CommandFactory, Parser};
use std::io::Write;
use std::sync::Arc;

/// Internal CLI state. Mirrors Go's `cli` struct.
pub struct Cli {
    pub format: String,
    pub indent: String,
    pub sort: bool,
    pub version: String,
    pub flags: PathFolderFlags,
    /// Failure-injection switch used by the `init_vaku_client` tests.
    pub fail: String,
    pub client: Option<Arc<dyn ClientInterface>>,
}

impl Cli {
    fn new(version: &str) -> Self {
        Self {
            format: "text".into(),
            indent: "    ".into(),
            sort: true,
            version: version.into(),
            flags: PathFolderFlags::default(),
            fail: String::new(),
            client: None,
        }
    }

    fn ctx(&self) -> OutputCtx<'_> {
        OutputCtx {
            format: &self.format,
            indent: &self.indent,
            sort: self.sort,
        }
    }
}

/// Run the CLI against a pre-parsed list of `args` (not including the binary
/// name) and the supplied writers. Returns the exit code.
pub fn run(version: &str, args: &[String], out: &mut dyn Write, err: &mut dyn Write) -> u8 {
    run_with_client(version, args, out, err, None)
}

/// Like [`run`] but lets tests inject a pre-built [`ClientInterface`] so
/// they don't need a real Vault server.
pub fn run_with_client(
    version: &str,
    args: &[String],
    out: &mut dyn Write,
    err: &mut dyn Write,
    client: Option<Arc<dyn ClientInterface>>,
) -> u8 {
    let mut argv = vec!["vaku".to_string()];
    argv.extend_from_slice(args);

    // Two-pass parse: first try, on error print clap-style help to the right
    // writer and return.
    let parsed = match VakuArgs::try_parse_from(&argv) {
        Ok(p) => p,
        Err(e) => {
            // Distinguish help/version from real errors so the exit code and
            // target writer match Go's cobra defaults.
            let kind = e.kind();
            use clap::error::ErrorKind::*;
            let exit = match kind {
                DisplayHelp | DisplayVersion | DisplayHelpOnMissingArgumentOrSubcommand => 0,
                _ => 1,
            };
            // Clap's `Error::print()` writes to stderr by default; reroute.
            let msg = e.render().to_string();
            if exit == 0 {
                let _ = out.write_all(msg.as_bytes());
            } else {
                let _ = err.write_all(msg.as_bytes());
            }
            return exit;
        }
    };

    let mut cli = Cli::new(version);
    cli.client = client;
    cli.format = parsed.format.clone();
    cli.indent = parsed.indent_char.clone();
    cli.sort = parsed.sort;

    let rt = tokio::runtime::Builder::new_multi_thread()
        .enable_all()
        .build()
        .expect("tokio runtime");
    rt.block_on(async move {
        match parsed.cmd {
            None => {
                // No subcommand — emulate cobra's default of showing help.
                let mut cmd = VakuArgs::command();
                let help = cmd.render_long_help().to_string();
                let _ = out.write_all(help.as_bytes());
                0
            }
            Some(TopCmd::Version) => run_version(&cli, out),
            Some(TopCmd::Docs { path }) => match docs::generate_markdown_tree(&path) {
                Ok(()) => 0,
                Err(_) => {
                    let ctx = cli.ctx();
                    output(&ctx, Out::Err(ERR_DOC_GEN_MARKDOWN.to_string()), out, err);
                    1
                }
            },
            Some(TopCmd::Completion { shell }) => {
                let mut cmd = VakuArgs::command();
                let bin_name = cmd.get_name().to_string();
                clap_complete::generate(shell, &mut cmd, bin_name, out);
                0
            }
            Some(TopCmd::Path(root)) => {
                cli.flags = root.flags;
                if let Err(code) = validate_flags(&cli, err) {
                    return code;
                }
                if cli.client.is_none() {
                    match init_client(&cli) {
                        Ok(c) => cli.client = Some(c),
                        Err(msg) => {
                            let _ = writeln!(err, "ERROR: {msg}");
                            return 1;
                        }
                    }
                }
                dispatch_path(&mut cli, root.cmd, out, err).await
            }
            Some(TopCmd::Folder(root)) => {
                cli.flags = root.flags;
                if let Err(code) = validate_flags(&cli, err) {
                    return code;
                }
                if cli.client.is_none() {
                    match init_client(&cli) {
                        Ok(c) => cli.client = Some(c),
                        Err(msg) => {
                            let _ = writeln!(err, "ERROR: {msg}");
                            return 1;
                        }
                    }
                }
                dispatch_folder(&mut cli, root.cmd, out, err).await
            }
        }
    })
}

fn run_version(cli: &Cli, out: &mut dyn Write) -> u8 {
    let ctx = cli.ctx();
    let mut m = serde_json::Map::new();
    m.insert(
        "API".into(),
        serde_json::Value::String(crate::api::version::version().into()),
    );
    m.insert("CLI".into(), serde_json::Value::String(cli.version.clone()));
    let mut dummy = std::io::sink();
    output(&ctx, Out::Map(m), out, &mut dummy);
    0
}

fn validate_flags(cli: &Cli, err: &mut dyn Write) -> Result<(), u8> {
    // format
    if cli.format != "text" && cli.format != "json" {
        let _ = writeln!(err, "ERROR: {}", ERR_FLAG_INVALID_FORMAT);
        return Err(1);
    }
    if cli.flags.workers < 1 {
        let _ = writeln!(err, "ERROR: {}", ERR_FLAG_INVALID_WORKERS);
        return Err(1);
    }

    // mount pair validation
    validate_mount_pair(
        &cli.flags.mount_path,
        &cli.flags.mount_version,
        "2",
        ERR_FLAG_MOUNT_VERSION_NO_PATH,
        ERR_FLAG_INVALID_MOUNT_VERSION,
        err,
    )?;
    validate_mount_pair(
        &cli.flags.mount_path_source,
        &cli.flags.mount_version_source,
        "2",
        ERR_FLAG_SRC_MOUNT_VERSION_NO_PATH,
        ERR_FLAG_INVALID_SRC_MOUNT_VERSION,
        err,
    )?;
    validate_mount_pair(
        &cli.flags.mount_path_destination,
        &cli.flags.mount_version_destination,
        "2",
        ERR_FLAG_DST_MOUNT_VERSION_NO_PATH,
        ERR_FLAG_INVALID_DST_MOUNT_VERSION,
        err,
    )?;
    Ok(())
}

fn validate_mount_pair(
    path: &str,
    version: &str,
    default_version: &str,
    err_no_path: CliErr,
    err_invalid_version: CliErr,
    w: &mut dyn Write,
) -> Result<(), u8> {
    if version != default_version && !version.is_empty() && path.is_empty() {
        let _ = writeln!(w, "ERROR: {err_no_path}");
        return Err(1);
    }
    if !path.is_empty() && !is_valid_mount_version(version) {
        let _ = writeln!(w, "ERROR: {err_invalid_version}");
        return Err(1);
    }
    Ok(())
}

fn is_valid_mount_version(v: &str) -> bool {
    v == "1" || v == "2"
}

fn init_client(cli: &Cli) -> Result<Arc<dyn ClientInterface>, String> {
    let mut builder = Client::builder()
        .with_absolute_path(cli.flags.absolute_path)
        .with_ignore_access_errors(cli.flags.ignore_read_errors)
        .with_workers(cli.flags.workers);

    let src_addr = effective_addr(&cli.flags.source_address, &cli.flags.address);
    let src_token = resolve_token(&cli.flags.source_token, &cli.flags.token, cli)?;
    let src_ns = effective_addr(&cli.flags.source_namespace, &cli.flags.namespace);

    let src_http = build_http(&src_addr, &src_token, &src_ns, cli)?;
    builder = builder.with_logical(Arc::new(src_http));

    if !cli.flags.destination_address.is_empty() || !cli.flags.destination_token.is_empty() {
        let dst_addr = &cli.flags.destination_address;
        let dst_token = resolve_token(&cli.flags.destination_token, "", cli)?;
        let dst_ns = &cli.flags.destination_namespace;
        let dst_http = build_http(dst_addr, &dst_token, dst_ns, cli)?;
        builder = builder.with_dst_logical(Arc::new(dst_http));
    }

    let src_mp = effective_addr(&cli.flags.mount_path_source, &cli.flags.mount_path);
    if !src_mp.is_empty() {
        let v = effective_addr(&cli.flags.mount_version_source, &cli.flags.mount_version);
        builder = builder.with_src_mount_provider(Arc::new(StaticMountProvider::new(&src_mp, &v)));
    }
    if !cli.flags.mount_path_destination.is_empty() {
        builder = builder.with_dst_mount_provider(Arc::new(StaticMountProvider::new(
            &cli.flags.mount_path_destination,
            &cli.flags.mount_version_destination,
        )));
    }

    let client = builder
        .build()
        .map_err(|e| format!("{ERR_INIT_VAKU_CLIENT}\n{}{e}", indent_for(cli)))?;
    Ok(Arc::new(client))
}

fn indent_for(cli: &Cli) -> &str {
    &cli.indent
}

fn effective_addr(specific: &str, alias: &str) -> String {
    if !specific.is_empty() {
        specific.to_string()
    } else {
        alias.to_string()
    }
}

fn build_http(
    addr: &str,
    token: &str,
    namespace: &str,
    cli: &Cli,
) -> Result<VaultHttpClient, String> {
    let env_addr = std::env::var("VAULT_ADDR").unwrap_or_default();
    let final_addr = if !addr.is_empty() {
        addr.to_string()
    } else if !env_addr.is_empty() {
        env_addr
    } else {
        "https://127.0.0.1:8200".to_string()
    };
    let final_ns = if !namespace.is_empty() {
        Some(namespace)
    } else {
        None
    };

    if cli.fail == "vault.NewClient" {
        return Err(format!(
            "{ERR_INIT_VAKU_CLIENT}\n{}{ERR_NEW_VAULT_CLIENT}",
            indent_for(cli)
        ));
    }

    let mut client = VaultHttpClient::new(&final_addr, token, final_ns).map_err(|e| {
        format!(
            "{ERR_INIT_VAKU_CLIENT}\n{}{ERR_SET_ADDRESS}\n{}{e}",
            indent_for(cli),
            indent_for(cli)
        )
    })?;
    // Allow overriding token from env or vault token file post-construction.
    if token.is_empty() {
        match resolve_default_token(cli) {
            Ok(t) if !t.is_empty() => {
                client = VaultHttpClient::new(&final_addr, &t, final_ns).map_err(|e| {
                    format!(
                        "{ERR_INIT_VAKU_CLIENT}\n{}{ERR_SET_ADDRESS}\n{}{e}",
                        indent_for(cli),
                        indent_for(cli)
                    )
                })?;
            }
            _ => {}
        }
    }
    Ok(client)
}

fn resolve_token(specific: &str, alias: &str, cli: &Cli) -> Result<String, String> {
    if !specific.is_empty() {
        return Ok(specific.to_string());
    }
    if !alias.is_empty() {
        return Ok(alias.to_string());
    }
    resolve_default_token(cli)
}

fn resolve_default_token(cli: &Cli) -> Result<String, String> {
    if cli.fail == "config.DefaultTokenHelper" {
        return Err(format!(
            "{ERR_INIT_VAKU_CLIENT}\n{}{ERR_SET_VAULT_TOKEN}\n{}{ERR_VAULT_TOKEN_HELPER}",
            indent_for(cli),
            indent_for(cli)
        ));
    }
    if let Ok(v) = std::env::var("VAULT_TOKEN") {
        if !v.is_empty() {
            return Ok(v);
        }
    }
    if cli.fail == "helper.Get" {
        return Err(format!(
            "{ERR_INIT_VAKU_CLIENT}\n{}{ERR_SET_VAULT_TOKEN}\n{}{ERR_GET_VAULT_TOKEN}",
            indent_for(cli),
            indent_for(cli)
        ));
    }
    // Fall back to ~/.vault-token, matching the Vault CLI default helper.
    if let Some(home) = dirs::home_dir() {
        let p = home.join(".vault-token");
        if let Ok(t) = std::fs::read_to_string(&p) {
            return Ok(t.trim().to_string());
        }
    }
    Ok(String::new())
}

async fn dispatch_path(
    cli: &mut Cli,
    sub: PathSub,
    out: &mut dyn Write,
    err: &mut dyn Write,
) -> u8 {
    let c = cli.client.clone().expect("client initialized");
    let ctx = cli.ctx();
    let mut sink = std::io::sink();
    match sub {
        PathSub::List { path } => match c.path_list(&path).await {
            Ok(list) => {
                output(&ctx, Out::List(list), out, &mut sink);
                0
            }
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        PathSub::Read { path } => match c.path_read(&path).await {
            Ok(Some(m)) => {
                output(&ctx, Out::Map(m), out, &mut sink);
                0
            }
            Ok(None) => 0,
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        PathSub::Delete { path } => match c.path_delete(&path).await {
            Ok(()) => 0,
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        PathSub::DeleteMeta { path } => match c.path_delete_meta(&path).await {
            Ok(()) => 0,
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        PathSub::Search { path, search } => match c.path_search(&path, &search).await {
            Ok(b) => {
                output(&ctx, Out::Text(b.to_string()), out, &mut sink);
                0
            }
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        PathSub::Copy {
            src,
            dst,
            all_versions,
        } => {
            let res = if all_versions {
                c.path_copy_all_versions(&src, &dst).await
            } else {
                c.path_copy(&src, &dst).await
            };
            match res {
                Ok(()) => 0,
                Err(e) => emit_err(&ctx, &e, out, err),
            }
        }
        PathSub::Move {
            src,
            dst,
            all_versions,
            destroy,
        } => {
            let res = if all_versions {
                c.path_move_all_versions(&src, &dst).await
            } else if destroy {
                match c.path_copy(&src, &dst).await {
                    Ok(()) => c.path_delete_meta(&src).await,
                    Err(e) => Err(e),
                }
            } else {
                c.path_move(&src, &dst).await
            };
            match res {
                Ok(()) => 0,
                Err(e) => emit_err(&ctx, &e, out, err),
            }
        }
        PathSub::Write { .. } | PathSub::Update { .. } | PathSub::Destroy { .. } => {
            // Hidden — mirror Go's "no-op help" behaviour.
            0
        }
    }
}

async fn dispatch_folder(
    cli: &mut Cli,
    sub: FolderSub,
    out: &mut dyn Write,
    err: &mut dyn Write,
) -> u8 {
    let c = cli.client.clone().expect("client initialized");
    let ctx = cli.ctx();
    let mut sink = std::io::sink();
    match sub {
        FolderSub::List { path } => match c.folder_list(&path).await {
            Ok(list) => {
                output(&ctx, Out::List(list), out, &mut sink);
                0
            }
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        FolderSub::Read { path } => match c.folder_read(&path).await {
            Ok(Some(m)) => {
                output(&ctx, Out::NestedMap(m), out, &mut sink);
                0
            }
            Ok(None) => 0,
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        FolderSub::Write { json } => {
            let parsed: Result<
                std::collections::BTreeMap<String, serde_json::Map<String, serde_json::Value>>,
                _,
            > = serde_json::from_str(&json);
            match parsed {
                Ok(input) => match c.folder_write(input).await {
                    Ok(()) => 0,
                    Err(e) => emit_err(&ctx, &e, out, err),
                },
                Err(je) => {
                    let combined = format!("{ERR_JSON_UNMARSHAL}\n{}{je}", cli.indent);
                    output(&ctx, Out::Err(combined), out, err);
                    1
                }
            }
        }
        FolderSub::Delete { path } => match c.folder_delete(&path).await {
            Ok(()) => 0,
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        FolderSub::DeleteMeta { path } => match c.folder_delete_meta(&path).await {
            Ok(()) => 0,
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        FolderSub::Search { path, search } => match c.folder_search(&path, &search).await {
            Ok(list) => {
                output(&ctx, Out::List(list), out, &mut sink);
                0
            }
            Err(e) => emit_err(&ctx, &e, out, err),
        },
        FolderSub::Copy {
            src,
            dst,
            all_versions,
        } => {
            let res = if all_versions {
                c.folder_copy_all_versions(&src, &dst).await
            } else {
                c.folder_copy(&src, &dst).await
            };
            match res {
                Ok(()) => 0,
                Err(e) => emit_err(&ctx, &e, out, err),
            }
        }
        FolderSub::Move {
            src,
            dst,
            all_versions,
            destroy,
        } => {
            let res = if all_versions {
                c.folder_move_all_versions(&src, &dst).await
            } else if destroy {
                match c.folder_copy(&src, &dst).await {
                    Ok(()) => c.folder_delete_meta(&src).await,
                    Err(e) => Err(e),
                }
            } else {
                c.folder_move(&src, &dst).await
            };
            match res {
                Ok(()) => 0,
                Err(e) => emit_err(&ctx, &e, out, err),
            }
        }
        FolderSub::Destroy { .. } => 0,
    }
}

fn emit_err(
    ctx: &OutputCtx<'_>,
    e: &dyn std::error::Error,
    out: &mut dyn Write,
    err: &mut dyn Write,
) -> u8 {
    output(ctx, Out::Err(e.to_string()), out, err);
    1
}

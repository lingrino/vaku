//! Seeded test paths, shared clients, and Vault container lifecycle.
//!
//! Mirrors Go's `api/main_test.go`. Two Vault dev containers (src + dst) are
//! launched once per test binary and reused. Each test allocates a unique
//! prefix number, seeds the canonical map into both servers at
//! `kv1/<N>/...` and `kv2/<N>/...`, and runs against those isolated paths.

use super::docker::VaultServer;
use super::injector::LogicalInjector;
use once_cell::sync::Lazy;
use serde_json::{json, Map, Value};
use std::collections::BTreeMap;
use std::sync::atomic::{AtomicUsize, Ordering};
use std::sync::Arc;
use vaku::api::client::Client;
use vaku::api::error::Error;
use vaku::api::logical::{Logical, VaultHttpClient};
use vaku::api::mount_provider::MountProvider;

/// Sentinel passed by tests that want to exercise the "no matching mount"
/// path. [`seeded_prefixes`] returns a single empty prefix for this.
pub const MOUNTLESS: &str = "mountless";

/// Canonical seed data. Same as Go's `seeds`.
fn seed_map() -> BTreeMap<&'static str, BTreeMap<&'static str, &'static str>> {
    BTreeMap::from_iter([
        ("0/1", BTreeMap::from_iter([("2", "3")])),
        ("0/4/5", BTreeMap::from_iter([("6", "7")])),
        ("0/4/8", BTreeMap::from_iter([("9", "10"), ("11", "12")])),
        ("0/4/13/14", BTreeMap::from_iter([("15", "16")])),
        (
            "0/4/13/17",
            BTreeMap::from_iter([("18", "19"), ("20", "21"), ("22", "23")]),
        ),
        ("0/4/13/24/25/26/27", BTreeMap::from_iter([("28", "29")])),
    ])
}

/// Holds two Vault dev containers — src + dst — for the duration of the
/// integration-test binary. `Lazy` ensures both are started exactly once.
pub static SERVERS: Lazy<TestServers> = Lazy::new(TestServers::start);

pub struct TestServers {
    pub src: Arc<VaultServer>,
    pub dst: Arc<VaultServer>,
}

impl TestServers {
    fn start() -> Self {
        // Start the two containers in parallel — they each take a few seconds.
        let src_h = std::thread::spawn(VaultServer::start);
        let dst_h = std::thread::spawn(VaultServer::start);
        let src = Arc::new(src_h.join().expect("src container"));
        let dst = Arc::new(dst_h.join().expect("dst container"));
        Self { src, dst }
    }
}

/// Monotonic counter for test prefixes. Starts at 100 to match Go.
static PATH_PREFIX: AtomicUsize = AtomicUsize::new(100);

/// Track whether KV mounts have been created on src/dst. They only need to be
/// created once per container.
static MOUNTS_READY: tokio::sync::OnceCell<()> = tokio::sync::OnceCell::const_new();

/// Bundled shared clients (mirrors Go's `sharedVaku` + `sharedVakuClean`).
pub struct SharedClients {
    /// Client with the path-marker injector spliced into the source `Logical`.
    /// Used by tests that need to inject errors/edge-case responses.
    pub vaku: Client,
    /// Same client wiring without the injector — used to seed and verify state.
    pub clean: Client,
}

/// Build a fresh pair of shared clients pointed at the test containers.
/// (Cheap; reqwest clients don't share state across our tests.)
pub async fn shared_clients() -> SharedClients {
    ensure_mounts().await;
    build_clients(false)
}

fn build_http_client(server: &VaultServer) -> Arc<dyn Logical> {
    Arc::new(VaultHttpClient::new(&server.addr, &server.token, None).expect("vault client"))
}

fn build_clients(_skip_injector: bool) -> SharedClients {
    let src_http = build_http_client(&SERVERS.src);
    let dst_http = build_http_client(&SERVERS.dst);

    let injector = Arc::new(LogicalInjector::new(src_http.clone(), false));
    let dst_injector = Arc::new(LogicalInjector::new(dst_http.clone(), false));

    let vaku = Client::builder()
        .with_logical(injector.clone())
        .with_dst_logical(dst_injector.clone())
        .with_absolute_path(false)
        .with_ignore_access_errors(false)
        .with_workers(5)
        .build()
        .expect("vaku");

    // Clean client uses raw HTTP — no injector.
    let clean = Client::builder()
        .with_logical(src_http)
        .with_dst_logical(dst_http)
        .with_absolute_path(false)
        .with_ignore_access_errors(false)
        .with_workers(5)
        .build()
        .expect("vaku clean");

    SharedClients { vaku, clean }
}

async fn ensure_mounts() {
    MOUNTS_READY
        .get_or_init(|| async {
            mount_both(&SERVERS.src).await;
            mount_both(&SERVERS.dst).await;
        })
        .await;
}

async fn mount_both(server: &VaultServer) {
    // Both kv1 and kv2 — idempotently. Vault returns 4xx with
    // "path is already in use" if we re-mount; tolerate that.
    for version in ["1", "2"] {
        let url = format!("{}/v1/sys/mounts/kv{version}/", server.addr);
        let body = json!({
            "type": "kv",
            "options": { "version": version },
        });
        let client = reqwest::Client::new();
        let resp = client
            .post(&url)
            .header("X-Vault-Token", &server.token)
            .json(&body)
            .send()
            .await
            .expect("mount request");
        let status = resp.status().as_u16();
        if !(status == 204 || status == 200 || status == 400) {
            let text = resp.text().await.unwrap_or_default();
            panic!("unexpected mount status {status} on kv{version}: {text}");
        }
    }
}

/// Allocate a unique test prefix and seed the canonical map into kv1+kv2 on
/// **both** the src and dst servers. Returns the per-version prefixes — e.g.
/// `vec!["kv1/101", "kv2/101"]` — that the caller can iterate over.
///
/// When `p == "mountless"`, returns a single empty prefix.
pub async fn seeded_prefixes(p: &str) -> Vec<String> {
    let clients = shared_clients().await;
    let prefix_num = PATH_PREFIX.fetch_add(1, Ordering::SeqCst);

    if p == MOUNTLESS {
        return vec![String::new()];
    }

    let mut prefixes = Vec::with_capacity(2);
    let seed = seed_map();
    for ver in ["1", "2"] {
        let mount_prefix = format!("kv{ver}/{prefix_num}");
        let mut write_map: BTreeMap<String, Map<String, Value>> = BTreeMap::new();
        for (rel, kvs) in &seed {
            let path = format!("{mount_prefix}/{rel}");
            let mut inner = Map::new();
            for (k, v) in kvs {
                inner.insert((*k).to_string(), json!(*v));
            }
            write_map.insert(path, inner);
        }

        clients
            .clean
            .folder_write(write_map.clone())
            .await
            .expect("seed src");
        clients
            .clean
            .as_destination()
            .folder_write(write_map)
            .await
            .expect("seed dst");

        prefixes.push(mount_prefix);
    }
    prefixes
}

/// Produces 4 (src_prefix, dst_prefix) pairs covering same-mount and
/// cross-mount combinations. Each pair operates on a freshly seeded prefix.
pub async fn seeded_prefix_product() -> [(String, String); 4] {
    let p1 = seeded_prefixes("").await;
    let p2 = seeded_prefixes("").await;
    let p3 = seeded_prefixes("").await;
    let p4 = seeded_prefixes("").await;
    [
        (p1[0].clone(), p1[0].clone()),
        (p2[0].clone(), p2[1].clone()),
        (p3[1].clone(), p3[0].clone()),
        (p4[1].clone(), p4[1].clone()),
    ]
}

/// Test name helper mirroring Go's `testName`.
pub fn test_name(sp: &str, dp: Option<&str>) -> String {
    match dp {
        Some(d) if !d.is_empty() => format!("~{sp}->{d}~"),
        _ => format!("~{sp}~"),
    }
}

/// Convenience helper for tests: a no-op `MountProvider` that always returns
/// `ErrNoMount` (matches Go's `mountless` path).
#[derive(Debug)]
pub struct NoMountProvider;

#[async_trait::async_trait]
impl MountProvider for NoMountProvider {
    async fn list_mounts(&self) -> Result<Vec<vaku::api::mount_provider::Mount>, Error> {
        Ok(Vec::new())
    }
}

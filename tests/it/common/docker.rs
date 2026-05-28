//! Docker-backed Vault server for integration tests.

use std::io::{Read, Write};
use std::net::TcpStream;
use std::process::{Command, Stdio};
use std::time::{Duration, Instant};

/// A running `hashicorp/vault:latest` container in dev mode. Removed on drop.
pub struct VaultServer {
    pub addr: String,
    pub token: String,
    pub port: u16,
    container_id: String,
}

impl VaultServer {
    /// Start a new Vault dev container. Polls `sys/health` until the server
    /// responds. Panics on failure — the caller (integration tests) treats
    /// these as test setup failures.
    pub fn start() -> Self {
        let token = "root".to_string();

        let output = Command::new("docker")
            .args([
                "run",
                "--rm",
                "-d",
                "-e",
                &format!("VAULT_DEV_ROOT_TOKEN_ID={token}"),
                "-e",
                "VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200",
                "-e",
                "SKIP_CHOWN=true",
                "-e",
                "SKIP_SETCAP=true",
                "--cap-add=IPC_LOCK",
                "-P",
                "hashicorp/vault:latest",
            ])
            .stdout(Stdio::piped())
            .stderr(Stdio::piped())
            .output()
            .expect("failed to spawn docker run");
        if !output.status.success() {
            panic!(
                "docker run failed: status={} stdout={} stderr={}",
                output.status,
                String::from_utf8_lossy(&output.stdout),
                String::from_utf8_lossy(&output.stderr),
            );
        }
        let container_id = String::from_utf8_lossy(&output.stdout).trim().to_string();

        let port = docker_port(&container_id, 8200);
        let addr = format!("http://127.0.0.1:{port}");
        wait_for_vault(&addr, port);

        Self {
            addr,
            token,
            port,
            container_id,
        }
    }
}

impl Drop for VaultServer {
    fn drop(&mut self) {
        let _ = Command::new("docker")
            .args(["rm", "-f", &self.container_id])
            .stdout(Stdio::null())
            .stderr(Stdio::null())
            .status();
    }
}

fn docker_port(container_id: &str, port: u16) -> u16 {
    let out = Command::new("docker")
        .args(["port", container_id, &format!("{port}/tcp")])
        .output()
        .expect("docker port");
    if !out.status.success() {
        panic!(
            "docker port failed: {}",
            String::from_utf8_lossy(&out.stderr)
        );
    }
    let text = String::from_utf8_lossy(&out.stdout);
    let line = text
        .lines()
        .find(|l| {
            l.starts_with("0.0.0.0") || l.starts_with("127.0.0.1")
        })
        .unwrap_or_else(|| panic!("no port mapping found in: {text}"));
    let port_str = line.rsplit(':').next().expect("port str");
    port_str
        .trim()
        .parse()
        .unwrap_or_else(|e| panic!("parse port '{port_str}': {e}"))
}

fn wait_for_vault(_addr: &str, port: u16) {
    let deadline = Instant::now() + Duration::from_secs(60);
    loop {
        if check_health(port) {
            return;
        }
        if Instant::now() >= deadline {
            panic!("vault did not become ready within 60s on port {port}");
        }
        std::thread::sleep(Duration::from_millis(200));
    }
}

/// Minimal sync HTTP/1.0 GET that succeeds when Vault returns any health
/// status code (200 sealed/standby/etc.). Avoids depending on reqwest in sync
/// context — keeps the test harness dependency-light.
fn check_health(port: u16) -> bool {
    let addr = format!("127.0.0.1:{port}");
    let Ok(mut stream) = TcpStream::connect_timeout(
        &addr.parse().expect("socket addr"),
        Duration::from_secs(1),
    ) else {
        return false;
    };
    let _ = stream.set_read_timeout(Some(Duration::from_secs(2)));
    let _ = stream.set_write_timeout(Some(Duration::from_secs(2)));
    let req = b"GET /v1/sys/health HTTP/1.0\r\nHost: 127.0.0.1\r\nConnection: close\r\n\r\n";
    if stream.write_all(req).is_err() {
        return false;
    }
    let mut buf = [0u8; 64];
    let n = stream.read(&mut buf).unwrap_or(0);
    let prefix = &buf[..n];
    // Look for "HTTP/1.0 " then a digit. Any 2xx/4xx response counts as ready
    // (Vault returns 200 unsealed, 429 standby, 472/473 etc.).
    prefix
        .windows(9)
        .any(|w| w.starts_with(b"HTTP/1.") && w[8].is_ascii_digit())
}

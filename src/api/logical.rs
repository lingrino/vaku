//! The Vault HTTP surface used by Vaku.
//!
//! [`Logical`] is the trait abstraction; [`VaultHttpClient`] is the concrete
//! `reqwest`-based implementation. Both mirror the (small) subset of the
//! Vault `Logical()` Go interface that Vaku actually uses.

use crate::api::error::{BoxError, Error, ErrorKind};
use async_trait::async_trait;
use serde_json::{Map, Value};
use std::fmt;
use url::Url;

/// A Vault secret: just the `data` field. Everything else in a Vault response
/// is ignored — Vaku doesn't need it.
#[derive(Debug, Clone, Default)]
pub struct Secret {
    pub data: Option<Map<String, Value>>,
}

/// The minimal Vault HTTP surface Vaku relies on.
#[async_trait]
pub trait Logical: Send + Sync {
    async fn read(&self, path: &str) -> Result<Option<Secret>, BoxError>;
    async fn read_with_data(
        &self,
        path: &str,
        params: &[(&str, &str)],
    ) -> Result<Option<Secret>, BoxError>;
    async fn list(&self, path: &str) -> Result<Option<Secret>, BoxError>;
    async fn write(&self, path: &str, data: Value) -> Result<Option<Secret>, BoxError>;
    async fn delete(&self, path: &str) -> Result<Option<Secret>, BoxError>;
}

/// HTTP error from a non-2xx Vault response.
#[derive(Debug)]
pub struct VaultHttpError {
    pub status: u16,
    pub body: String,
}

impl fmt::Display for VaultHttpError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        if self.body.is_empty() {
            write!(f, "vault returned status {}", self.status)
        } else {
            write!(f, "vault status {}: {}", self.status, self.body.trim())
        }
    }
}

impl std::error::Error for VaultHttpError {}

/// Concrete `reqwest`-based Vault client.
#[derive(Debug, Clone)]
pub struct VaultHttpClient {
    base: Url,
    token: String,
    namespace: Option<String>,
    client: reqwest::Client,
}

impl VaultHttpClient {
    /// Create a new client. `addr` should be a full URL like
    /// `http://127.0.0.1:8200`.
    pub fn new(addr: &str, token: &str, namespace: Option<&str>) -> Result<Self, Error> {
        let mut base = Url::parse(addr).map_err(|e| {
            Error::wrap(
                "vault address",
                ErrorKind::Custom("invalid vault address".into()),
                Some(Box::new(e)),
            )
        })?;

        // Ensure a trailing slash on the base URL so `base.join` works.
        if !base.path().ends_with('/') {
            let mut path = base.path().to_owned();
            path.push('/');
            base.set_path(&path);
        }

        let client = reqwest::Client::builder().build().map_err(|e| {
            Error::wrap(
                "vault client",
                ErrorKind::Custom("http client build".into()),
                Some(Box::new(e)),
            )
        })?;

        Ok(Self {
            base,
            token: token.to_string(),
            namespace: namespace.map(str::to_string),
            client,
        })
    }

    fn url(&self, path: &str) -> Result<Url, BoxError> {
        let trimmed = path.trim_start_matches('/');
        let combined = format!("v1/{trimmed}");
        Ok(self.base.join(&combined)?)
    }

    fn req(&self, method: reqwest::Method, url: Url) -> reqwest::RequestBuilder {
        let mut rb = self.client.request(method, url);
        if !self.token.is_empty() {
            rb = rb.header("X-Vault-Token", &self.token);
        }
        if let Some(ns) = &self.namespace {
            if !ns.is_empty() {
                rb = rb.header("X-Vault-Namespace", ns);
            }
        }
        rb
    }

    async fn send(&self, req: reqwest::RequestBuilder) -> Result<Option<Secret>, BoxError> {
        let resp = req.send().await?;
        let status = resp.status();
        if status.as_u16() == 404 {
            // Vault returns 404 for "secret deleted" and "secret missing".
            // Match Go's `vault/api`: a missing secret yields `(nil, nil)`.
            return Ok(None);
        }
        if !status.is_success() {
            let body = resp.text().await.unwrap_or_default();
            return Err(Box::new(VaultHttpError {
                status: status.as_u16(),
                body,
            }));
        }

        let body = resp.text().await?;
        if body.is_empty() {
            return Ok(None);
        }
        // Vault always wraps responses as `{"data": {...}, ...}`. We only care
        // about the `data` key for non-listing reads; for LIST, `data` contains
        // a `keys` array. Either way we surface the `data` map verbatim.
        let parsed: Value = serde_json::from_str(&body)?;
        let data = parsed.get("data").and_then(Value::as_object).cloned();
        Ok(Some(Secret { data }))
    }
}

#[async_trait]
impl Logical for VaultHttpClient {
    async fn read(&self, path: &str) -> Result<Option<Secret>, BoxError> {
        let url = self.url(path)?;
        self.send(self.req(reqwest::Method::GET, url)).await
    }

    async fn read_with_data(
        &self,
        path: &str,
        params: &[(&str, &str)],
    ) -> Result<Option<Secret>, BoxError> {
        let mut url = self.url(path)?;
        if !params.is_empty() {
            let mut q = url.query_pairs_mut();
            for (k, v) in params {
                q.append_pair(k, v);
            }
        }
        self.send(self.req(reqwest::Method::GET, url)).await
    }

    async fn list(&self, path: &str) -> Result<Option<Secret>, BoxError> {
        let mut url = self.url(path)?;
        url.query_pairs_mut().append_pair("list", "true");
        self.send(self.req(reqwest::Method::GET, url)).await
    }

    async fn write(&self, path: &str, data: Value) -> Result<Option<Secret>, BoxError> {
        let url = self.url(path)?;
        let req = self.req(reqwest::Method::PUT, url).json(&data);
        self.send(req).await
    }

    async fn delete(&self, path: &str) -> Result<Option<Secret>, BoxError> {
        let url = self.url(path)?;
        self.send(self.req(reqwest::Method::DELETE, url)).await
    }
}

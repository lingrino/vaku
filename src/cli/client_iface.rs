//! Trait that the CLI talks to. Mirror of Go's `ClientInterface`.
//!
//! Production code uses [`vaku::Client`] which implements it. Tests provide
//! a mock implementation so CLI behaviour can be verified without spinning
//! up Vault.

use crate::api::client::Client;
use crate::api::error::Error;
use crate::api::secret::SecretMeta;
use async_trait::async_trait;
use serde_json::{Map, Value};
use std::collections::BTreeMap;

#[async_trait]
pub trait ClientInterface: Send + Sync {
    async fn path_list(&self, p: &str) -> Result<Vec<String>, Error>;
    async fn path_read(&self, p: &str) -> Result<Option<Map<String, Value>>, Error>;
    async fn path_read_meta(&self, p: &str) -> Result<Option<SecretMeta>, Error>;
    async fn path_read_version(&self, p: &str, v: i64)
        -> Result<Option<Map<String, Value>>, Error>;
    async fn path_write(&self, p: &str, data: Option<Map<String, Value>>) -> Result<(), Error>;
    async fn path_delete(&self, p: &str) -> Result<(), Error>;
    async fn path_delete_meta(&self, p: &str) -> Result<(), Error>;
    async fn path_destroy(&self, p: &str, versions: &[i64]) -> Result<(), Error>;
    async fn path_update(&self, p: &str, data: Option<Map<String, Value>>) -> Result<(), Error>;
    async fn path_search(&self, p: &str, s: &str) -> Result<bool, Error>;
    async fn path_copy(&self, src: &str, dst: &str) -> Result<(), Error>;
    async fn path_copy_all_versions(&self, src: &str, dst: &str) -> Result<(), Error>;
    async fn path_move(&self, src: &str, dst: &str) -> Result<(), Error>;
    async fn path_move_all_versions(&self, src: &str, dst: &str) -> Result<(), Error>;

    async fn folder_list(&self, p: &str) -> Result<Vec<String>, Error>;
    async fn folder_read(
        &self,
        p: &str,
    ) -> Result<Option<BTreeMap<String, Map<String, Value>>>, Error>;
    async fn folder_write(&self, data: BTreeMap<String, Map<String, Value>>) -> Result<(), Error>;
    async fn folder_delete(&self, p: &str) -> Result<(), Error>;
    async fn folder_delete_meta(&self, p: &str) -> Result<(), Error>;
    async fn folder_destroy(&self, p: &str, versions: &[i64]) -> Result<(), Error>;
    async fn folder_search(&self, p: &str, s: &str) -> Result<Vec<String>, Error>;
    async fn folder_copy(&self, src: &str, dst: &str) -> Result<(), Error>;
    async fn folder_copy_all_versions(&self, src: &str, dst: &str) -> Result<(), Error>;
    async fn folder_move(&self, src: &str, dst: &str) -> Result<(), Error>;
    async fn folder_move_all_versions(&self, src: &str, dst: &str) -> Result<(), Error>;
}

#[async_trait]
impl ClientInterface for Client {
    async fn path_list(&self, p: &str) -> Result<Vec<String>, Error> {
        Client::path_list(self, p).await
    }
    async fn path_read(&self, p: &str) -> Result<Option<Map<String, Value>>, Error> {
        Client::path_read(self, p).await
    }
    async fn path_read_meta(&self, p: &str) -> Result<Option<SecretMeta>, Error> {
        Client::path_read_meta(self, p).await
    }
    async fn path_read_version(
        &self,
        p: &str,
        v: i64,
    ) -> Result<Option<Map<String, Value>>, Error> {
        Client::path_read_version(self, p, v).await
    }
    async fn path_write(&self, p: &str, data: Option<Map<String, Value>>) -> Result<(), Error> {
        Client::path_write(self, p, data).await
    }
    async fn path_delete(&self, p: &str) -> Result<(), Error> {
        Client::path_delete(self, p).await
    }
    async fn path_delete_meta(&self, p: &str) -> Result<(), Error> {
        Client::path_delete_meta(self, p).await
    }
    async fn path_destroy(&self, p: &str, versions: &[i64]) -> Result<(), Error> {
        Client::path_destroy(self, p, versions).await
    }
    async fn path_update(&self, p: &str, data: Option<Map<String, Value>>) -> Result<(), Error> {
        Client::path_update(self, p, data).await
    }
    async fn path_search(&self, p: &str, s: &str) -> Result<bool, Error> {
        Client::path_search(self, p, s).await
    }
    async fn path_copy(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::path_copy(self, src, dst).await
    }
    async fn path_copy_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::path_copy_all_versions(self, src, dst).await
    }
    async fn path_move(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::path_move(self, src, dst).await
    }
    async fn path_move_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::path_move_all_versions(self, src, dst).await
    }

    async fn folder_list(&self, p: &str) -> Result<Vec<String>, Error> {
        Client::folder_list(self, p).await
    }
    async fn folder_read(
        &self,
        p: &str,
    ) -> Result<Option<BTreeMap<String, Map<String, Value>>>, Error> {
        Client::folder_read(self, p).await
    }
    async fn folder_write(&self, data: BTreeMap<String, Map<String, Value>>) -> Result<(), Error> {
        Client::folder_write(self, data).await
    }
    async fn folder_delete(&self, p: &str) -> Result<(), Error> {
        Client::folder_delete(self, p).await
    }
    async fn folder_delete_meta(&self, p: &str) -> Result<(), Error> {
        Client::folder_delete_meta(self, p).await
    }
    async fn folder_destroy(&self, p: &str, versions: &[i64]) -> Result<(), Error> {
        Client::folder_destroy(self, p, versions).await
    }
    async fn folder_search(&self, p: &str, s: &str) -> Result<Vec<String>, Error> {
        Client::folder_search(self, p, s).await
    }
    async fn folder_copy(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::folder_copy(self, src, dst).await
    }
    async fn folder_copy_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::folder_copy_all_versions(self, src, dst).await
    }
    async fn folder_move(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::folder_move(self, src, dst).await
    }
    async fn folder_move_all_versions(&self, src: &str, dst: &str) -> Result<(), Error> {
        Client::folder_move_all_versions(self, src, dst).await
    }
}

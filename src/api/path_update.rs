//! `PathUpdate` — read-merge-write convenience.

use crate::api::client::Client;
use crate::api::error::{Error, ErrorKind};
use serde_json::{Map, Value};

impl Client {
    /// Update a path by reading existing data, merging `data` into it (new
    /// keys win), and writing back.
    pub async fn path_update(
        &self,
        p: &str,
        data: Option<Map<String, Value>>,
    ) -> Result<(), Error> {
        let Some(data) = data else {
            return Err(Error::wrap(
                p,
                ErrorKind::PathUpdate,
                Some(Box::new(Error::wrap("", ErrorKind::NilData, None))),
            ));
        };

        let mut read = match self.path_read(p).await {
            Ok(Some(m)) => m,
            Ok(None) => Map::new(),
            Err(e) => return Err(Error::wrap(p, ErrorKind::PathUpdate, Some(Box::new(e)))),
        };

        for (k, v) in data {
            read.insert(k, v);
        }

        self.path_write(p, Some(read))
            .await
            .map_err(|e| Error::wrap(p, ErrorKind::PathUpdate, Some(Box::new(e))))
    }
}

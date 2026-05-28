//! CLI-side error sentinels. Mirror Go's `errInit*`, etc.

use std::fmt;

#[derive(Debug)]
pub struct CliErr {
    pub msg: &'static str,
}

impl CliErr {
    pub const fn new(msg: &'static str) -> Self {
        Self { msg }
    }
}

impl fmt::Display for CliErr {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        f.write_str(self.msg)
    }
}

impl std::error::Error for CliErr {}

pub const ERR_INIT_VAKU_CLIENT: CliErr = CliErr::new("initializing vaku client");
pub const ERR_NEW_VAULT_CLIENT: CliErr = CliErr::new("creating new vault client");
pub const ERR_VAULT_TOKEN_HELPER: CliErr = CliErr::new("getting default token helper");
pub const ERR_GET_VAULT_TOKEN: CliErr = CliErr::new("using helper to get vault token");
pub const ERR_SET_VAULT_TOKEN: CliErr = CliErr::new("setting vault token");
pub const ERR_SET_ADDRESS: CliErr = CliErr::new("setting vault address");

pub const ERR_FLAG_INVALID_FORMAT: CliErr = CliErr::new("format must be one of: text|json");
pub const ERR_FLAG_INVALID_WORKERS: CliErr = CliErr::new("workers must be >= 1");
pub const ERR_FLAG_INVALID_MOUNT_VERSION: CliErr = CliErr::new("mount-version must be one of: 1|2");
pub const ERR_FLAG_MOUNT_VERSION_NO_PATH: CliErr =
    CliErr::new("mount-version requires --mount-path");
pub const ERR_FLAG_INVALID_SRC_MOUNT_VERSION: CliErr =
    CliErr::new("mount-version-source must be one of: 1|2");
pub const ERR_FLAG_SRC_MOUNT_VERSION_NO_PATH: CliErr =
    CliErr::new("mount-version-source requires --mount-path-source");
pub const ERR_FLAG_INVALID_DST_MOUNT_VERSION: CliErr =
    CliErr::new("mount-version-destination must be one of: 1|2");
pub const ERR_FLAG_DST_MOUNT_VERSION_NO_PATH: CliErr =
    CliErr::new("mount-version-destination requires --mount-path-destination");

pub const ERR_OUTPUT_FORMAT: CliErr = CliErr::new("unsupported output format");
pub const ERR_OUTPUT_TYPE: CliErr = CliErr::new("unsupported output type");
pub const ERR_JSON_MARSHAL: CliErr = CliErr::new("json marshal");
pub const ERR_JSON_UNMARSHAL: CliErr = CliErr::new("json unmarshal");
pub const ERR_DOC_GEN_MARKDOWN: CliErr = CliErr::new("failed to generate markdown docs");

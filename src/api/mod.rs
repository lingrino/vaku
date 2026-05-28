//! Vaku API — the public Rust library surface.
//!
//! Modules mirror the original Go layout, one file per operation.

pub mod client;
pub mod error;
pub mod folder_copy;
pub mod folder_copy_all_versions;
pub mod folder_delete;
pub mod folder_delete_meta;
pub mod folder_destroy;
pub mod folder_list;
pub mod folder_move;
pub mod folder_move_all_versions;
pub mod folder_read;
pub mod folder_search;
pub mod folder_write;
pub mod helpers;
pub mod logical;
pub mod mount_provider;
pub mod mounts;
pub mod path_copy;
pub mod path_copy_all_versions;
pub mod path_delete;
pub mod path_delete_meta;
pub mod path_destroy;
pub mod path_list;
pub mod path_move;
pub mod path_move_all_versions;
pub mod path_read;
pub mod path_read_meta;
pub mod path_search;
pub mod path_update;
pub mod path_write;
pub mod secret;
pub mod version;

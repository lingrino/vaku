//! Single integration-test binary so all `tests/it/*.rs` files share one set
//! of Vault dev containers (`SERVERS` in `common::seeds`).

pub mod common;

// Pure-function tests (no Vault container required).
mod api_error;
mod api_helpers;
mod api_version;

// Vault-backed tests. Compiled into the `it` binary; skipped at runtime if
// docker isn't available.
mod api_client;
mod api_folder_copy;
mod api_folder_copy_all_versions;
mod api_folder_delete;
mod api_folder_delete_meta;
mod api_folder_destroy;
mod api_folder_list;
mod api_folder_move;
mod api_folder_move_all_versions;
mod api_folder_read;
mod api_folder_search;
mod api_folder_write;
mod api_mounts;
mod api_path_copy;
mod api_path_copy_all_versions;
mod api_path_delete;
mod api_path_delete_meta;
mod api_path_destroy;
mod api_path_list;
mod api_path_move;
mod api_path_move_all_versions;
mod api_path_read;
mod api_path_read_meta;
mod api_path_search;
mod api_path_update;
mod api_path_write;

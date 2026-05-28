//! Smoke tests covering top-level CLI behaviour: `vaku`, `vaku path`,
//! `vaku folder`, `vaku --help`, and a few error paths.

use crate::cli_common::{run_cli, run_cli_with_mock};

#[test]
fn test_vaku_root_help() {
    let (_code, out, _err) = run_cli(&[]);
    // We expect the long description text to appear in output.
    assert!(out.contains("Vaku is a CLI for working with large Vault k/v secret engines"));
}

#[test]
fn test_vaku_unknown_subcommand_fails() {
    let (code, _out, _err) = run_cli(&["INVALID"]);
    assert_eq!(code, 1);
}

#[test]
fn test_vaku_path_help() {
    let (_code, out, _err) = run_cli(&["path", "--help"]);
    assert!(out.contains("Commands that act on Vault paths"));
}

#[test]
fn test_vaku_folder_help() {
    let (_code, out, _err) = run_cli(&["folder", "--help"]);
    assert!(out.contains("Commands that act on Vault folders"));
}

#[test]
fn test_path_list_mock() {
    let (code, out, err) = run_cli_with_mock(&["path", "list", "foo"]);
    assert_eq!(code, 0, "stderr: {err}");
    assert_eq!(out, "foo\nmoo\n");
}

#[test]
fn test_path_read_mock() {
    let (code, out, err) = run_cli_with_mock(&["path", "read", "foo"]);
    assert_eq!(code, 0, "stderr: {err}");
    assert_eq!(out, "biz: baz\nfoo: bar\n");
}

#[test]
fn test_path_delete_mock() {
    let (code, out, err) = run_cli_with_mock(&["path", "delete", "foo"]);
    assert_eq!(code, 0, "stderr: {err}");
    assert_eq!(out, "");
}

#[test]
fn test_path_delete_meta_mock() {
    let (code, out, err) = run_cli_with_mock(&["path", "delete-meta", "foo"]);
    assert_eq!(code, 0, "stderr: {err}");
    assert_eq!(out, "");
}

#[test]
fn test_path_search_mock() {
    let (code, out, _err) = run_cli_with_mock(&["path", "search", "foo", "bar"]);
    assert_eq!(code, 0);
    assert_eq!(out, "true\n");
}

#[test]
fn test_path_copy_mock() {
    let (code, out, _err) = run_cli_with_mock(&["path", "copy", "foo", "bar"]);
    assert_eq!(code, 0);
    assert_eq!(out, "");
    let (code, _, _) = run_cli_with_mock(&["path", "copy", "--all-versions", "foo", "bar"]);
    assert_eq!(code, 0);
}

#[test]
fn test_path_move_mock() {
    let (code, _, _) = run_cli_with_mock(&["path", "move", "foo", "bar"]);
    assert_eq!(code, 0);
    let (code, _, _) = run_cli_with_mock(&["path", "move", "--all-versions", "foo", "bar"]);
    assert_eq!(code, 0);
    let (code, _, _) = run_cli_with_mock(&["path", "move", "--destroy", "foo", "bar"]);
    assert_eq!(code, 0);
}

#[test]
fn test_folder_list_mock() {
    let (code, out, _) = run_cli_with_mock(&["folder", "list", "foo"]);
    assert_eq!(code, 0);
    assert_eq!(out, "bim/bom\nfoo/bar\nfoo/baz\n");
}

#[test]
fn test_folder_read_mock() {
    let (code, out, _) = run_cli_with_mock(&["folder", "read", "foo"]);
    assert_eq!(code, 0);
    assert_eq!(out, "bar\nhoo: boo\nfoo\nbim: bom\nbiz: baz\n");
}

#[test]
fn test_folder_write_mock() {
    let (code, out, _) = run_cli_with_mock(&["folder", "write", "{}"]);
    assert_eq!(code, 0);
    assert_eq!(out, "");

    let (code, out, _) = run_cli_with_mock(&["folder", "write", "{\"a/b/c\": {\"foo\": \"bar\"}}"]);
    assert_eq!(code, 0);
    assert_eq!(out, "");

    let (code, _out, err) = run_cli_with_mock(&["folder", "write", "error"]);
    assert_eq!(code, 1);
    assert!(err.starts_with("ERROR: json unmarshal\n"), "err was: {err}");
}

#[test]
fn test_folder_search_mock() {
    let (code, out, _) = run_cli_with_mock(&["folder", "search", "foo", "bar"]);
    assert_eq!(code, 0);
    assert_eq!(out, "bim/bom\nfoo/bar\n");
}

#[test]
fn test_folder_copy_mock() {
    let (code, _, _) = run_cli_with_mock(&["folder", "copy", "foo", "bar"]);
    assert_eq!(code, 0);
    let (code, _, _) = run_cli_with_mock(&["folder", "copy", "--all-versions", "foo", "bar"]);
    assert_eq!(code, 0);
}

#[test]
fn test_folder_move_mock() {
    let (code, _, _) = run_cli_with_mock(&["folder", "move", "foo", "bar"]);
    assert_eq!(code, 0);
    let (code, _, _) = run_cli_with_mock(&["folder", "move", "--all-versions", "foo", "bar"]);
    assert_eq!(code, 0);
    let (code, _, _) = run_cli_with_mock(&["folder", "move", "--destroy", "foo", "bar"]);
    assert_eq!(code, 0);
}

#[test]
fn test_version_subcommand() {
    let (code, out, _) = run_cli(&["version"]);
    assert_eq!(code, 0);
    assert_eq!(out, "API: 3.0.0\nCLI: dev\n");
}

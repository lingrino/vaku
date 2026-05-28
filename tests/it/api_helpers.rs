//! Ports `api/helpers_test.go` — pure-function path helpers, no Vault needed.

use serde_json::{json, Map, Value};
use std::collections::BTreeMap;
use vaku::api::helpers::{
    add_prefix, add_prefix_list, ensure_folder, ensure_prefix, ensure_prefix_list,
    ensure_prefix_map, insert_into_path, is_folder, path_join, trim_prefix_list, trim_prefix_map,
};

#[test]
fn test_path_join() {
    let cases: &[(&[&str], &str)] = &[
        (&["/"], "/"),
        (&["a/"], "a/"),
        (&["b", ""], "b"),
        (&["a/b", "c"], "a/b/c"),
        (&["d/e/", "/f"], "d/e/f"),
        (&["/g/h/", "/i/"], "g/h/i/"),
        (&["/j/", "/k/l", "m"], "j/k/l/m"),
    ];
    for (input, want) in cases {
        assert_eq!(path_join(input), *want, "input={input:?}");
    }
}

#[test]
fn test_is_folder() {
    let cases = [
        ("/", true),
        ("a/", true),
        ("a/b/", true),
        ("", false),
        ("a", false),
        ("a/b", false),
        ("123/456", false),
    ];
    for (input, want) in cases {
        assert_eq!(is_folder(input), want, "input={input}");
    }
}

#[test]
fn test_ensure_folder() {
    let cases = [("", "/"), ("a", "a/"), ("a/", "a/"), ("a/b", "a/b/")];
    for (input, want) in cases {
        assert_eq!(ensure_folder(input), want, "input={input}");
    }
}

#[test]
fn test_add_prefix() {
    let cases = [
        ("", "", ""),
        ("a", "", "a"),
        ("", "a", "a"),
        ("a/", "a", "a/a/"),
        ("a", "a/", "a/a"),
        ("a/b/c/d", "a/b/", "a/b/a/b/c/d"),
        ("a/b/c/d", "b", "b/a/b/c/d"),
    ];
    for (give, prefix, want) in cases {
        assert_eq!(
            add_prefix(give, prefix),
            want,
            "give={give} prefix={prefix}"
        );
    }
}

#[test]
fn test_ensure_prefix() {
    let cases = [
        ("", "", ""),
        ("a", "", "a"),
        ("", "a", "a"),
        ("a/", "a", "a/"),
        ("a", "a/", "a/a"),
        ("a/b/c/d", "a/b/", "a/b/c/d"),
        ("a/b/c/d", "b", "b/a/b/c/d"),
    ];
    for (give, prefix, want) in cases {
        assert_eq!(ensure_prefix(give, prefix), want);
    }
}

#[test]
fn test_add_prefix_list() {
    let cases: &[(&[&str], &str, &[&str])] = &[
        (&["a"], "b", &["b/a"]),
        (&["/c/d/e/"], "/f/", &["f/c/d/e/"]),
        (&["/g/"], "h", &["h/g/"]),
        (&["i/j"], "i", &["i/i/j"]),
    ];
    for (input, prefix, want) in cases {
        let mut got: Vec<String> = input.iter().map(|s| s.to_string()).collect();
        add_prefix_list(&mut got, prefix);
        let want_vec: Vec<String> = want.iter().map(|s| s.to_string()).collect();
        assert_eq!(got, want_vec, "prefix={prefix}");
    }
}

#[test]
fn test_ensure_prefix_list() {
    let cases: &[(&[&str], &str, &[&str])] = &[
        (&["a"], "b", &["b/a"]),
        (&["/c/d/e/"], "/f/", &["f/c/d/e/"]),
        (&["/g/"], "h", &["h/g/"]),
        (&["i/j"], "i", &["i/j"]),
    ];
    for (input, prefix, want) in cases {
        let mut got: Vec<String> = input.iter().map(|s| s.to_string()).collect();
        ensure_prefix_list(&mut got, prefix);
        let want_vec: Vec<String> = want.iter().map(|s| s.to_string()).collect();
        assert_eq!(got, want_vec, "prefix={prefix}");
    }
}

#[test]
fn test_trim_prefix_list() {
    let cases: &[(&[&str], &str, &[&str])] = &[
        (&["a"], "b", &["a"]),
        (&["/c/d/e/"], "/c/", &["d/e/"]),
        (&["f/g"], "f", &["g"]),
        (&["i/j"], "k", &["i/j"]),
    ];
    for (input, prefix, want) in cases {
        let mut got: Vec<String> = input.iter().map(|s| s.to_string()).collect();
        trim_prefix_list(&mut got, prefix);
        let want_vec: Vec<String> = want.iter().map(|s| s.to_string()).collect();
        assert_eq!(got, want_vec, "prefix={prefix}");
    }
}

fn mk_map(kvs: &[(&str, &str, &str)]) -> BTreeMap<String, Map<String, Value>> {
    let mut m: BTreeMap<String, Map<String, Value>> = BTreeMap::new();
    for (path, k, v) in kvs {
        let mut inner = Map::new();
        inner.insert((*k).to_string(), json!(*v));
        m.insert((*path).to_string(), inner);
    }
    m
}

#[test]
fn test_ensure_prefix_map() {
    let cases = [
        (
            "foo",
            vec![("foo/bar", "a", "b")],
            vec![("foo/bar", "a", "b")],
        ),
        (
            "foo/",
            vec![("foo/bar", "a", "b")],
            vec![("foo/bar", "a", "b")],
        ),
        (
            "fo",
            vec![("foo/bar", "a", "b")],
            vec![("foo/bar", "a", "b")],
        ),
        (
            "fooo",
            vec![("foo/bar", "a", "b")],
            vec![("fooo/foo/bar", "a", "b")],
        ),
    ];
    for (prefix, give, want) in cases {
        let mut got = mk_map(&give);
        let want_map = mk_map(&want);
        ensure_prefix_map(&mut got, prefix);
        assert_eq!(got, want_map, "prefix={prefix}");
    }
}

#[test]
fn test_trim_prefix_map() {
    let cases = [
        ("foo", vec![("foo/bar", "a", "b")], vec![("bar", "a", "b")]),
        ("foo/", vec![("foo/bar", "a", "b")], vec![("bar", "a", "b")]),
        ("fo", vec![("foo/bar", "a", "b")], vec![("o/bar", "a", "b")]),
        (
            "fooo",
            vec![("foo/bar", "a", "b")],
            vec![("foo/bar", "a", "b")],
        ),
    ];
    for (prefix, give, want) in cases {
        let mut got = mk_map(&give);
        let want_map = mk_map(&want);
        trim_prefix_map(&mut got, prefix);
        assert_eq!(got, want_map, "prefix={prefix}");
    }
}

#[test]
fn test_insert_into_path() {
    let cases = [
        ("", "", "", ""),
        ("foo", "fo", "b", "fo/b/o"),
        ("foo/bar", "fo", "b", "fo/b/o/bar"),
        ("foo/bar", "foo", "baz", "foo/baz/bar"),
        ("foo/bar/", "foo/", "baz/", "foo/baz/bar/"),
        ("1/2/3/4/5/6", "1/2/3", "foo", "1/2/3/foo/4/5/6"),
    ];
    for (path, after, insert, want) in cases {
        assert_eq!(
            insert_into_path(path, after, insert),
            want,
            "path={path} after={after} insert={insert}"
        );
    }
}

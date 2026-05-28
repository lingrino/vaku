//! Tests for the CLI helpers — output renderers and combine_err.

use serde_json::{json, Map, Value};
use std::collections::BTreeMap;
use vaku::cli::helpers::{combine_err, output, Out, OutputCtx};

fn run(format: &str, indent: &str, sort: bool, o: Out) -> (String, String) {
    let mut out: Vec<u8> = Vec::new();
    let mut err: Vec<u8> = Vec::new();
    let ctx = OutputCtx {
        format,
        indent,
        sort,
    };
    output(&ctx, o, &mut out, &mut err);
    (
        String::from_utf8(out).unwrap(),
        String::from_utf8(err).unwrap(),
    )
}

#[test]
fn output_text_nil() {
    let (o, e) = run("text", "", true, Out::None);
    assert_eq!(o, "");
    assert_eq!(e, "");
}

#[test]
fn output_text_string() {
    let (o, e) = run("text", "", true, Out::Text("foo".into()));
    assert_eq!(o, "foo\n");
    assert_eq!(e, "");
}

#[test]
fn output_text_list() {
    let (o, _) = run(
        "text",
        "",
        true,
        Out::List(vec!["foo".into(), "bar".into()]),
    );
    assert_eq!(o, "bar\nfoo\n");
}

#[test]
fn output_text_map() {
    let mut m = Map::new();
    m.insert("foo".into(), json!("fooValue"));
    m.insert("bar".into(), json!(100));
    let (o, _) = run("text", "", true, Out::Map(m));
    assert_eq!(o, "bar: 100\nfoo: fooValue\n");
}

#[test]
fn output_text_nested_map() {
    let mut foo = Map::new();
    foo.insert("infoo".into(), json!("fooValue"));
    foo.insert("inbar".into(), json!(100));
    let mut bar = Map::new();
    bar.insert("hello".into(), json!("world"));

    let mut outer: BTreeMap<String, Map<String, Value>> = BTreeMap::new();
    outer.insert("foo".into(), foo);
    outer.insert("bar".into(), bar);

    let (o, _) = run("text", "", true, Out::NestedMap(outer));
    assert_eq!(o, "bar\nhello: world\nfoo\ninbar: 100\ninfoo: fooValue\n");
}

#[test]
fn output_json_string() {
    let (o, _) = run("json", "", true, Out::Text("foo".into()));
    assert_eq!(o, "\"foo\"\n");
}

#[test]
fn output_json_list() {
    let (o, _) = run(
        "json",
        "",
        true,
        Out::List(vec!["foo".into(), "bar".into()]),
    );
    assert_eq!(o, "[\n\"foo\",\n\"bar\"\n]\n");
}

#[test]
fn output_json_map() {
    let mut m = Map::new();
    m.insert("foo".into(), json!("fooValue"));
    m.insert("bar".into(), json!(100));
    let (o, _) = run("json", "", true, Out::Map(m));
    assert_eq!(o, "{\n\"bar\": 100,\n\"foo\": \"fooValue\"\n}\n");
}

#[test]
fn output_json_nested_map() {
    let mut foo = Map::new();
    foo.insert("infoo".into(), json!("fooValue"));
    foo.insert("inbar".into(), json!(100));
    let mut bar = Map::new();
    bar.insert("hello".into(), json!("world"));

    let mut outer: BTreeMap<String, Map<String, Value>> = BTreeMap::new();
    outer.insert("foo".into(), foo);
    outer.insert("bar".into(), bar);

    let (o, _) = run("json", "", true, Out::NestedMap(outer));
    assert_eq!(
        o,
        "{\n\"bar\": {\n\"hello\": \"world\"\n},\n\"foo\": {\n\"inbar\": 100,\n\"infoo\": \"fooValue\"\n}\n}\n"
    );
}

#[test]
fn output_text_error() {
    let (o, e) = run("text", "", true, Out::Err("test error".into()));
    assert_eq!(o, "");
    assert_eq!(e, "ERROR: test error\n");
}

#[test]
fn output_json_error() {
    let (o, e) = run("json", "", true, Out::Err("test error".into()));
    assert_eq!(o, "");
    assert_eq!(e, "{\n\"error\": \"test error\"\n}\n");
}

#[test]
fn output_bad_format() {
    let (o, e) = run("invalid", "", true, Out::Text("".into()));
    assert_eq!(o, "");
    assert_eq!(e, "ERROR: unsupported output format\n");
}

#[test]
fn combine_err_table() {
    assert_eq!(combine_err(None, None, ""), None);
    assert_eq!(combine_err(Some("foo"), None, ""), Some("foo".into()));
    assert_eq!(combine_err(None, Some("bar"), ""), Some("bar".into()));
    assert_eq!(
        combine_err(Some("foo"), Some("bar"), ""),
        Some("foo\nbar".into())
    );
}

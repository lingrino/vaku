//! Output rendering helpers — text and JSON formatters.
//!
//! Reproduces Go's `cmd/helpers.go` output rules. For JSON we serialize with
//! a custom `serde_json::ser::Formatter` so the byte output matches Go's
//! `json.MarshalIndent("", indent)`: keys sorted, `": "` between key/value,
//! newlines between elements (or `[]`/`{}` for empty containers).

use crate::cli::errors::{ERR_JSON_MARSHAL, ERR_OUTPUT_FORMAT};
use serde::Serialize;
use serde_json::ser::{Formatter, Serializer};
use serde_json::{Map, Value};
use std::collections::BTreeMap;
use std::io::Write;

/// Anything the CLI may print.
pub enum Out {
    None,
    Text(String),
    List(Vec<String>),
    Map(Map<String, Value>),
    NestedMap(BTreeMap<String, Map<String, Value>>),
    Err(String),
}

pub struct OutputCtx<'a> {
    pub format: &'a str,
    pub indent: &'a str,
    pub sort: bool,
}

/// Combine two errors with newline+indent. Mirrors Go's `combineErr`.
pub fn combine_err(e1: Option<&str>, e2: Option<&str>, indent: &str) -> Option<String> {
    match (e1, e2) {
        (None, None) => None,
        (Some(a), None) => Some(a.to_string()),
        (None, Some(b)) => Some(b.to_string()),
        (Some(a), Some(b)) => Some(format!("{a}\n{indent}{b}")),
    }
}

pub fn output(ctx: &OutputCtx<'_>, out: Out, w: &mut dyn Write, e: &mut dyn Write) {
    match ctx.format {
        "json" => output_json(ctx, out, w, e),
        "text" => output_text(ctx, out, w, e),
        _ => write_err(e, &ERR_OUTPUT_FORMAT.to_string()),
    }
}

fn output_text(ctx: &OutputCtx<'_>, out: Out, w: &mut dyn Write, e: &mut dyn Write) {
    match out {
        Out::None => {}
        Out::Text(s) => {
            let _ = writeln!(w, "{s}");
        }
        Out::List(mut l) => {
            if ctx.sort {
                l.sort();
            }
            for s in l {
                let _ = writeln!(w, "{s}");
            }
        }
        Out::Map(m) => print_map_text(w, 0, ctx, &m),
        Out::NestedMap(m) => {
            let mut keys: Vec<&String> = m.keys().collect();
            if ctx.sort {
                keys.sort();
            }
            for k in keys {
                let _ = writeln!(w, "{k}");
                print_map_text(w, 1, ctx, &m[k]);
            }
        }
        Out::Err(msg) => write_err(e, &msg),
    }
}

fn print_map_text(
    w: &mut dyn Write,
    indent_times: usize,
    ctx: &OutputCtx<'_>,
    m: &Map<String, Value>,
) {
    let indent = ctx.indent.repeat(indent_times);
    let mut keys: Vec<&String> = m.keys().collect();
    if ctx.sort {
        keys.sort();
    }
    for k in keys {
        let v = &m[k];
        let _ = writeln!(w, "{indent}{k}: {}", display_value(v));
    }
}

fn display_value(v: &Value) -> String {
    match v {
        Value::String(s) => s.clone(),
        Value::Number(n) => n.to_string(),
        Value::Bool(b) => b.to_string(),
        Value::Null => "<nil>".to_string(),
        Value::Array(_) | Value::Object(_) => v.to_string(),
    }
}

fn output_json(ctx: &OutputCtx<'_>, out: Out, w: &mut dyn Write, e: &mut dyn Write) {
    let res = match out {
        Out::None => return,
        Out::Err(msg) => {
            let mut payload = Map::new();
            payload.insert("error".into(), Value::String(msg));
            write_json(e, ctx.indent, &Value::Object(payload))
        }
        Out::Text(s) => write_json(w, ctx.indent, &Value::String(s)),
        Out::List(l) => {
            let arr = Value::Array(l.into_iter().map(Value::String).collect());
            write_json(w, ctx.indent, &arr)
        }
        Out::Map(m) => write_json(w, ctx.indent, &Value::Object(m)),
        Out::NestedMap(m) => {
            let mut outer = Map::new();
            for (k, inner) in m {
                outer.insert(k, Value::Object(inner));
            }
            write_json(w, ctx.indent, &Value::Object(outer))
        }
    };
    if res.is_err() {
        write_err(e, &ERR_JSON_MARSHAL.to_string());
    }
}

/// Serialize `v` byte-for-byte like Go's `json.MarshalIndent("", indent)`.
fn write_json(w: &mut dyn Write, indent: &str, v: &Value) -> std::io::Result<()> {
    // Object keys must be sorted to match Go. `serde_json::Value::Object`
    // built from BTreeMap is already sorted, but Map<String, Value> (with
    // preserve_order feature) preserves insertion order. We sort explicitly.
    let sorted = sort_value(v);

    let formatter = GoMarshalIndent::default().with_indent(indent);
    let mut ser = Serializer::with_formatter(WriteAdapter(w), formatter);
    sorted.serialize(&mut ser).map_err(std::io::Error::other)?;
    let WriteAdapter(w) = ser.into_inner();
    w.write_all(b"\n")?;
    Ok(())
}

fn sort_value(v: &Value) -> Value {
    match v {
        Value::Object(m) => {
            let mut sorted: Vec<(String, Value)> =
                m.iter().map(|(k, v)| (k.clone(), sort_value(v))).collect();
            sorted.sort_by(|a, b| a.0.cmp(&b.0));
            let mut out = Map::new();
            for (k, v) in sorted {
                out.insert(k, v);
            }
            Value::Object(out)
        }
        Value::Array(a) => Value::Array(a.iter().map(sort_value).collect()),
        _ => v.clone(),
    }
}

fn write_err(e: &mut dyn Write, msg: &str) {
    let _ = writeln!(e, "ERROR: {msg}");
}

/// Custom JSON pretty-formatter mimicking Go's `encoding/json.MarshalIndent`.
#[derive(Default)]
struct GoMarshalIndent<'a> {
    indent: &'a str,
    depth: usize,
    has_items: Vec<bool>,
}

impl<'a> GoMarshalIndent<'a> {
    fn with_indent(mut self, indent: &'a str) -> Self {
        self.indent = indent;
        self
    }

    fn newline_indent<W: ?Sized + std::io::Write>(&self, w: &mut W) -> std::io::Result<()> {
        w.write_all(b"\n")?;
        for _ in 0..self.depth {
            w.write_all(self.indent.as_bytes())?;
        }
        Ok(())
    }
}

impl Formatter for GoMarshalIndent<'_> {
    fn begin_array<W: ?Sized + std::io::Write>(&mut self, w: &mut W) -> std::io::Result<()> {
        self.depth += 1;
        self.has_items.push(false);
        w.write_all(b"[")
    }
    fn end_array<W: ?Sized + std::io::Write>(&mut self, w: &mut W) -> std::io::Result<()> {
        let had = self.has_items.pop().unwrap_or(false);
        self.depth -= 1;
        if had {
            self.newline_indent(w)?;
        }
        w.write_all(b"]")
    }
    fn begin_array_value<W: ?Sized + std::io::Write>(
        &mut self,
        w: &mut W,
        first: bool,
    ) -> std::io::Result<()> {
        if !first {
            w.write_all(b",")?;
        }
        if let Some(last) = self.has_items.last_mut() {
            *last = true;
        }
        self.newline_indent(w)
    }
    fn end_array_value<W: ?Sized + std::io::Write>(&mut self, _w: &mut W) -> std::io::Result<()> {
        Ok(())
    }

    fn begin_object<W: ?Sized + std::io::Write>(&mut self, w: &mut W) -> std::io::Result<()> {
        self.depth += 1;
        self.has_items.push(false);
        w.write_all(b"{")
    }
    fn end_object<W: ?Sized + std::io::Write>(&mut self, w: &mut W) -> std::io::Result<()> {
        let had = self.has_items.pop().unwrap_or(false);
        self.depth -= 1;
        if had {
            self.newline_indent(w)?;
        }
        w.write_all(b"}")
    }
    fn begin_object_key<W: ?Sized + std::io::Write>(
        &mut self,
        w: &mut W,
        first: bool,
    ) -> std::io::Result<()> {
        if !first {
            w.write_all(b",")?;
        }
        if let Some(last) = self.has_items.last_mut() {
            *last = true;
        }
        self.newline_indent(w)
    }
    fn begin_object_value<W: ?Sized + std::io::Write>(&mut self, w: &mut W) -> std::io::Result<()> {
        w.write_all(b": ")
    }
}

/// Adapter so `Serializer::with_formatter` accepts a `&mut dyn Write`.
struct WriteAdapter<'a>(&'a mut dyn Write);
impl std::io::Write for WriteAdapter<'_> {
    fn write(&mut self, buf: &[u8]) -> std::io::Result<usize> {
        self.0.write(buf)
    }
    fn flush(&mut self) -> std::io::Result<()> {
        self.0.flush()
    }
}

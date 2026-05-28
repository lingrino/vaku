//! Markdown documentation emitter — produces one `.md` per (sub)command
//! under the given directory. Output style mirrors cobra's
//! `doc.GenMarkdownTree`: short, synopsis, usage block, examples, options,
//! inherited options, and SEE ALSO links.

use crate::cli::args::VakuArgs;
use crate::cli::errors::ERR_DOC_GEN_MARKDOWN;
use clap::{Arg, ArgAction, Command, CommandFactory};
use std::fs;
use std::io::Write as _;
use std::path::{Path, PathBuf};

pub fn generate_markdown_tree(dir: &Path) -> Result<(), &'static str> {
    fs::create_dir_all(dir).map_err(|_| ERR_DOC_GEN_MARKDOWN.msg)?;
    let cmd = VakuArgs::command();
    walk_and_emit(&cmd, dir, &[]).map_err(|_| ERR_DOC_GEN_MARKDOWN.msg)
}

fn walk_and_emit(cmd: &Command, dir: &Path, parents: &[&str]) -> std::io::Result<()> {
    // File name: vaku.md, vaku_path.md, vaku_path_list.md ...
    let name_chain: Vec<&str> = parents
        .iter()
        .copied()
        .chain(std::iter::once(cmd.get_name()))
        .collect();
    let file_stem = name_chain.join("_");
    let path = dir.join(format!("{file_stem}.md"));

    let mut s = String::new();
    write_markdown(&mut s, cmd, &name_chain);

    let mut f = fs::File::create(&path)?;
    f.write_all(s.as_bytes())?;

    for sub in cmd.get_subcommands() {
        if sub.is_hide_set() {
            continue;
        }
        let mut new_parents: Vec<&str> = parents.to_vec();
        new_parents.push(cmd.get_name());
        walk_and_emit(sub, dir, &new_parents)?;
    }
    Ok(())
}

fn write_markdown(s: &mut String, cmd: &Command, name_chain: &[&str]) {
    let full = name_chain.join(" ");
    s.push_str(&format!("## {full}\n\n"));

    let about = cmd
        .get_about()
        .map(|a| a.to_string())
        .unwrap_or_else(|| full.clone());
    s.push_str(&format!("{about}\n\n"));

    s.push_str("### Synopsis\n\n");
    let long = cmd
        .get_long_about()
        .map(|a| a.to_string())
        .unwrap_or_else(|| about.clone());
    s.push_str(&format!("{long}\n\n"));

    if cmd.get_subcommands().count() == 0 {
        s.push_str("```\n");
        s.push_str(&usage_line(cmd, name_chain));
        s.push_str("\n```\n\n");
    }

    if let Some(ex) = cmd.get_after_help().map(|a| a.to_string()) {
        if !ex.is_empty() {
            s.push_str("### Examples\n\n```\n");
            s.push_str(&ex);
            s.push_str("\n```\n\n");
        }
    }

    let (local, inherited) = split_args(cmd);
    if !local.is_empty() {
        s.push_str("### Options\n\n```\n");
        for a in &local {
            s.push_str(&render_arg(a));
            s.push('\n');
        }
        s.push_str("```\n\n");
    }
    if !inherited.is_empty() {
        s.push_str("### Options inherited from parent commands\n\n```\n");
        for a in &inherited {
            s.push_str(&render_arg(a));
            s.push('\n');
        }
        s.push_str("```\n\n");
    }

    // SEE ALSO
    let mut see_also: Vec<String> = Vec::new();
    if name_chain.len() > 1 {
        let parent = name_chain[..name_chain.len() - 1].join("_");
        let parent_pretty = name_chain[..name_chain.len() - 1].join(" ");
        see_also.push(format!(
            "* [{parent_pretty}]({parent}.md)\t - {}",
            "Vaku is a CLI for working with large Vault k/v secret engines"
        ));
    }
    for sub in cmd.get_subcommands() {
        if sub.is_hide_set() {
            continue;
        }
        let sub_chain: Vec<&str> = name_chain
            .iter()
            .copied()
            .chain(std::iter::once(sub.get_name()))
            .collect();
        let sub_file = sub_chain.join("_");
        let sub_pretty = sub_chain.join(" ");
        let about = sub.get_about().map(|a| a.to_string()).unwrap_or_default();
        see_also.push(format!("* [{sub_pretty}]({sub_file}.md)\t - {about}"));
    }
    if !see_also.is_empty() {
        s.push_str("### SEE ALSO\n\n");
        for l in &see_also {
            s.push_str(l);
            s.push('\n');
        }
        s.push('\n');
    }
}

fn usage_line(cmd: &Command, name_chain: &[&str]) -> String {
    let positionals: Vec<String> = cmd
        .get_positionals()
        .filter(|a| !a.is_hide_set())
        .map(|a| format!("<{}>", a.get_id().as_str()))
        .collect();
    let prefix = name_chain.join(" ");
    if positionals.is_empty() {
        format!("{prefix} [flags]")
    } else {
        format!("{prefix} {} [flags]", positionals.join(" "))
    }
}

fn split_args(cmd: &Command) -> (Vec<&Arg>, Vec<&Arg>) {
    let mut local = Vec::new();
    let mut inherited = Vec::new();
    for a in cmd.get_arguments() {
        if a.is_positional() || a.is_hide_set() {
            continue;
        }
        if a.is_global_set() {
            inherited.push(a);
        } else {
            local.push(a);
        }
    }
    local.sort_by_key(|a| a.get_id().as_str().to_string());
    inherited.sort_by_key(|a| a.get_id().as_str().to_string());
    (local, inherited)
}

fn render_arg(a: &Arg) -> String {
    let mut parts: Vec<String> = Vec::new();
    if let Some(short) = a.get_short() {
        parts.push(format!("-{short}"));
    } else {
        parts.push("  ".to_string());
    }
    let long = a.get_long().unwrap_or_else(|| a.get_id().as_str());
    parts.push(format!("--{long}"));

    let help = a.get_help().map(|h| h.to_string()).unwrap_or_default();
    let value_hint = match a.get_action() {
        ArgAction::SetTrue | ArgAction::SetFalse => String::new(),
        _ => " string".to_string(),
    };
    let default = a
        .get_default_values()
        .first()
        .map(|s| s.to_string_lossy().to_string())
        .filter(|s| !s.is_empty())
        .map(|s| format!(" (default \"{s}\")"))
        .unwrap_or_default();
    format!("  {}{}   {help}{}", parts.join(", "), value_hint, default)
}

#[allow(dead_code)]
fn _suppress() -> PathBuf {
    PathBuf::new()
}

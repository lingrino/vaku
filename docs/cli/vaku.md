## vaku

Vaku is a CLI for working with large Vault k/v secret engines

### Synopsis

Vaku is a CLI for working with large Vault k/v secret engines

The Vaku CLI provides path- and folder-based commands that work on
both Version 1 and Version 2 K/V secret engines. Vaku can help manage
large amounts of Vault data by updating secrets in place, moving
paths or folders, searching secrets, and more.

Vaku is not a replacement for the Vault CLI and requires that you are
already authenticated to Vault before running any commands. Vaku
commands should not be run on non-K/V engines.

CLI documentation - 'vaku help [cmd]'
API documentation - https://docs.rs/vaku
Built by Sean Lingren <sean@lingren.com>

### Examples

```
vaku folder list secret/foo
```

### Options inherited from parent commands

```
    , --format string   output format: text|json (default "text")
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort string   sort output text (default "true")
```

### SEE ALSO

* [vaku path](vaku_path.md)	 - Commands that act on Vault paths
* [vaku folder](vaku_folder.md)	 - Commands that act on Vault folders
* [vaku version](vaku_version.md)	 - Print vaku version


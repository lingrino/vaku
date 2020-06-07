## vaku

Vaku is a CLI for working with large Vault k/v secret engines

### Synopsis

Vaku is a CLI for working with large Vault k/v secret engines

The Vaku CLI provides path and folder based commands that work on
both version 1 and version 2 k/v secret engines. Vaku can help manage
large amounts of vault data by updating secrets in place, moving
paths or folders, searching secrets, and more.

Vaku is not a replacement for the Vault CLI and requires that you
already are authenticated to Vault before running any commands. Vaku
commands should not be run on non-k/v engines.

CLI documentation - 'vaku help [cmd]'
API documentation - https://pkg.go.dev/github.com/lingrino/vaku/vaku
Built by Sean Lingren <sean@lingrino.com>

### Examples

```
vaku folder list secret/foo
```

### Options

```
      --format string        output format: text|json (default "text")
  -h, --help                 help for vaku
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort                 sort output text (default true)
```

### SEE ALSO

* [vaku completion](vaku_completion.md)	 - Generates shell completions
* [vaku path](vaku_path.md)	 - Commands that act on Vault paths
* [vaku version](vaku_version.md)	 - Print vaku version


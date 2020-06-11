## vaku path

Commands that act on Vault paths

### Synopsis

Commands that act on Vault paths

Commands under the path subcommand act on Vault paths. Vaku can list,
copy, move, search, etc.. on Vault paths.

### Examples

```
vaku path list secret/foo
```

### Options

```
  -p, --absolute-path                show absolute path in output
  -a, --address string               address of the Vault server
      --destination-address string   address of the destination Vault server
      --destination-token string     token for the destination vault server (alias for --token)
  -h, --help                         help for path
      --source-address string        address of the source Vault server (alias for --address)
      --source-token string          token for the source vault server (alias for --token)
  -t, --token string                 token for the vault server
  -w, --workers int                  number of concurrent workers (default 10)
```

### Options inherited from parent commands

```
      --format string        output format: text|json (default "text")
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort                 sort output text (default true)
```

### SEE ALSO

* [vaku](vaku.md)	 - Vaku is a CLI for working with large Vault k/v secret engines
* [vaku path delete](vaku_path_delete.md)	 - Delete all paths at a path
* [vaku path list](vaku_path_list.md)	 - List all paths at a path
* [vaku path read](vaku_path_read.md)	 - Read all paths at a path


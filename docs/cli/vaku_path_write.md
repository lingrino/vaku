## vaku path write

Write all paths at a path

### Synopsis

Write all paths at a path

```
vaku path write <path> [flags]
```

### Examples

```
vaku path write secret/foo
```

### Options

```
  -h, --help   help for write
```

### Options inherited from parent commands

```
  -p, --absolute-path                show absolute path in output
  -a, --address string               address of the Vault server
      --destination-address string   address of the destination Vault server
      --destination-token string     token for the destination vault server (alias for --token)
      --format string                output format: text|json (default "text")
  -i, --indent-char string           string used for indents (default "    ")
  -s, --sort                         sort output text (default true)
      --source-address string        address of the source Vault server (alias for --address)
      --source-token string          token for the source vault server (alias for --token)
  -t, --token string                 token for the vault server
  -w, --workers int                  number of concurrent workers (default 10)
```

### SEE ALSO

* [vaku path](vaku_path.md)	 - Commands that act on Vault paths


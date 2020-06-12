## vaku folder read

Recursively read all paths in a folder

### Synopsis

Recursively read all paths in a folder

```
vaku folder read <path> [flags]
```

### Examples

```
vaku folder read secret/foo
```

### Options

```
  -h, --help   help for read
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

* [vaku folder](vaku_folder.md)	 - Commands that act on Vault folders


## vaku path copy

Copy a secret from a source path to a destination path

### Synopsis

Copy a secret from a source path to a destination path

```
vaku path copy <source path> <destination path> [flags]
```

### Examples

```
vaku path copy secret/foo secret/bar
```

### Options

```
      --all-versions   copy all versions of the secret (KV v2 only)
  -h, --help           help for copy
```

### Options inherited from parent commands

```
  -p, --absolute-path                  show absolute path in output
  -a, --address string                 address of the Vault server
      --destination-address string     address of the destination Vault server
      --destination-namespace string   name of the vault namespace to use in the destination client
      --destination-token string       token for the destination vault server (alias for --token)
      --format string                  output format: text|json (default "text")
      --ignore-read-errors             ignore path read errors and continue
  -i, --indent-char string             string used for indents (default "    ")
  -m, --mount-path string              mount path to use (bypasses sys/mounts lookup)
      --mount-version string           mount version: 1|2 (requires --mount-path) (default "2")
  -n, --namespace string               name of the vault namespace to use in the source client
  -s, --sort                           sort output text (default true)
      --source-address string          address of the source Vault server (alias for --address)
      --source-namespace string        name of the vault namespace to use in the source client (alias for --namespace)
      --source-token string            token for the source vault server (alias for --token)
  -t, --token string                   token for the vault server
  -w, --workers int                    number of concurrent workers (default 10)
```

### SEE ALSO

* [vaku path](vaku_path.md)	 - Commands that act on Vault paths


## vaku path read

Read a secret at a path

### Synopsis

Read a secret at a path

```
vaku path read <path> [flags]
```

### Examples

```
vaku path read secret/foo
```

### Options

```
  -h, --help   help for read
```

### Options inherited from parent commands

```
  -p, --absolute-path                      show absolute path in output
  -a, --address string                     address of the Vault server
      --destination-address string         address of the destination Vault server
      --destination-namespace string       name of the vault namespace to use in the destination client
      --destination-token string           token for the destination vault server (alias for --token)
      --format string                      output format: text|json (default "text")
      --ignore-read-errors                 ignore path read errors and continue
  -i, --indent-char string                 string used for indents (default "    ")
  -m, --mount-path string                  source mount path (bypasses sys/mounts lookup, alias for --mount-path-source)
      --mount-path-destination string      destination mount path (bypasses sys/mounts lookup)
      --mount-path-source string           source mount path (bypasses sys/mounts lookup)
      --mount-version string               source mount version: 1|2 (requires --mount-path, alias for --mount-version-source) (default "2")
      --mount-version-destination string   destination mount version: 1|2 (requires --mount-path-destination) (default "2")
      --mount-version-source string        source mount version: 1|2 (requires --mount-path-source) (default "2")
  -n, --namespace string                   name of the vault namespace to use in the source client
  -s, --sort                               sort output text (default true)
      --source-address string              address of the source Vault server (alias for --address)
      --source-namespace string            name of the vault namespace to use in the source client (alias for --namespace)
      --source-token string                token for the source vault server (alias for --token)
  -t, --token string                       token for the vault server
  -w, --workers int                        number of concurrent workers (default 10)
```

### SEE ALSO

* [vaku path](vaku_path.md)	 - Commands that act on Vault paths


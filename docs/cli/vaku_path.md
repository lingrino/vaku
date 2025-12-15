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
  -p, --absolute-path                      show absolute path in output
  -a, --address string                     address of the Vault server
      --destination-address string         address of the destination Vault server
      --destination-namespace string       name of the vault namespace to use in the destination client
      --destination-token string           token for the destination vault server (alias for --token)
  -h, --help                               help for path
      --ignore-read-errors                 ignore path read errors and continue
  -m, --mount-path string                  source mount path (bypasses sys/mounts lookup, alias for --mount-path-source)
      --mount-path-destination string      destination mount path (bypasses sys/mounts lookup)
      --mount-path-source string           source mount path (bypasses sys/mounts lookup)
      --mount-version string               source mount version: 1|2 (requires --mount-path, alias for --mount-version-source) (default "2")
      --mount-version-destination string   destination mount version: 1|2 (requires --mount-path-destination) (default "2")
      --mount-version-source string        source mount version: 1|2 (requires --mount-path-source) (default "2")
  -n, --namespace string                   name of the vault namespace to use in the source client
      --source-address string              address of the source Vault server (alias for --address)
      --source-namespace string            name of the vault namespace to use in the source client (alias for --namespace)
      --source-token string                token for the source vault server (alias for --token)
  -t, --token string                       token for the vault server
  -w, --workers int                        number of concurrent workers (default 10)
```

### Options inherited from parent commands

```
      --format string        output format: text|json (default "text")
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort                 sort output text (default true)
```

### SEE ALSO

* [vaku](vaku.md)	 - Vaku is a CLI for working with large Vault k/v secret engines
* [vaku path copy](vaku_path_copy.md)	 - Copy a secret from a source path to a destination path
* [vaku path delete](vaku_path_delete.md)	 - Delete a secret at a path
* [vaku path delete-meta](vaku_path_delete-meta.md)	 - Delete all secret metadata and versions at a path. V2 engines only.
* [vaku path list](vaku_path_list.md)	 - List all paths at a path
* [vaku path move](vaku_path_move.md)	 - Move a secret from a source path to a destination path
* [vaku path read](vaku_path_read.md)	 - Read a secret at a path
* [vaku path search](vaku_path_search.md)	 - Search a secret for a search string


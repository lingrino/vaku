## vaku folder

Commands that act on Vault folders

### Synopsis

Commands that act on Vault folders

Commands under the folder subcommand act on Vault folders. Folders
are designated by paths that end in a '/' such as 'secret/foo/'. Vaku
can list, copy, move, search, etc.. on Vault folders.

### Examples

```
vaku folder list secret/foo
```

### Options

```
  -p, --absolute-path                  show absolute path in output
  -a, --address string                 address of the Vault server
      --destination-address string     address of the destination Vault server
      --destination-namespace string   name of the vault namespace to use in the destination client
      --destination-token string       token for the destination vault server (alias for --token)
  -h, --help                           help for folder
      --ignore-read-errors             ignore path read errors and continue
  -n, --namespace string               name of the vault namespace to use in the source client
      --source-address string          address of the source Vault server (alias for --address)
      --source-namespace string        name of the vault namespace to use in the source client (alias for --namespace)
      --source-token string            token for the source vault server (alias for --token)
  -t, --token string                   token for the vault server
  -w, --workers int                    number of concurrent workers (default 10)
```

### Options inherited from parent commands

```
      --format string        output format: text|json (default "text")
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort                 sort output text (default true)
```

### SEE ALSO

* [vaku](vaku.md)	 - Vaku is a CLI for working with large Vault k/v secret engines
* [vaku folder copy](vaku_folder_copy.md)	 - Recursively copy all secrets in source folder to destination folder
* [vaku folder delete](vaku_folder_delete.md)	 - Recursively delete all secrets in a folder
* [vaku folder delete-meta](vaku_folder_delete-meta.md)	 - Recursively delete all secrets metadata and versions in a folder. V2 engines only.
* [vaku folder list](vaku_folder_list.md)	 - Recursively list all paths in a folder
* [vaku folder move](vaku_folder_move.md)	 - Recursively move all secrets in source folder to destination folder
* [vaku folder read](vaku_folder_read.md)	 - Recursively read all secrets in a folder
* [vaku folder search](vaku_folder_search.md)	 - Recursively search all secrets in a folder for a search string
* [vaku folder write](vaku_folder_write.md)	 - write a folder of secrets. WARNING: command expects a very specific json input


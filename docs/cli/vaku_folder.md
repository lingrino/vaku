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
  -p, --absolute-path                show absolute path in output
  -a, --address string               address of the Vault server
      --destination-address string   address of the destination Vault server
      --destination-token string     token for the destination vault server (alias for --token)
  -h, --help                         help for folder
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
* [vaku folder copy](vaku_folder_copy.md)	 - Copy a folder from source to destination
* [vaku folder delete](vaku_folder_delete.md)	 - Recursively delete all paths in a folder
* [vaku folder list](vaku_folder_list.md)	 - Recursively list all paths at a path
* [vaku folder move](vaku_folder_move.md)	 - Move a folder from source to destination
* [vaku folder read](vaku_folder_read.md)	 - Recursively read all paths in a folder
* [vaku folder search](vaku_folder_search.md)	 - Search for a secret in a folder


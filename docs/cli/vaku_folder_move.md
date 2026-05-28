## vaku folder move

Recursively move all secrets in source folder to destination folder

### Synopsis

Recursively move all secrets in source folder to destination folder

```
vaku folder move <src> <dst> [flags]
```

### Examples

```
vaku folder move secret/foo secret/bar
```

### Options

```
    , --all-versions   move all versions of the secret (KV v2 only)
    , --destroy   permanently destroy all versions at source after copy (KV v2 only)
```

### SEE ALSO

* [vaku folder](vaku_folder.md)	 - Vaku is a CLI for working with large Vault k/v secret engines


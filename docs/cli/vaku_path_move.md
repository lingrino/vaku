## vaku path move

Move a secret from a source path to a destination path

### Synopsis

Move a secret from a source path to a destination path

```
vaku path move <src> <dst> [flags]
```

### Examples

```
vaku path move secret/foo secret/bar
```

### Options

```
    , --all-versions   move all versions of the secret (KV v2 only)
    , --destroy   permanently destroy all versions at source after copy (KV v2 only)
```

### SEE ALSO

* [vaku path](vaku_path.md)	 - Vaku is a CLI for working with large Vault k/v secret engines


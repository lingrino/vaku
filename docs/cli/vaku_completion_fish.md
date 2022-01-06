## vaku completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	vaku completion fish | source

To load completions for every new session, execute once:

	vaku completion fish > ~/.config/fish/completions/vaku.fish

You will need to start a new shell for this setup to take effect.


```
vaku completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --format string        output format: text|json (default "text")
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort                 sort output text (default true)
```

### SEE ALSO

* [vaku completion](vaku_completion.md)	 - Generate the autocompletion script for the specified shell


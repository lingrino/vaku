## vaku completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(vaku completion zsh); compdef _vaku vaku

To load completions for every new session, execute once:

#### Linux:

	vaku completion zsh > "${fpath[1]}/_vaku"

#### macOS:

	vaku completion zsh > $(brew --prefix)/share/zsh/site-functions/_vaku

You will need to start a new shell for this setup to take effect.


```
vaku completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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


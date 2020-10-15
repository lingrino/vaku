## vaku completion

Generate shell completions

### Synopsis

To install completions for your shell

# Bash: In ~/.bashrc
source <(vaku completion bash)

# Fish: In ~/.config/fish/config.fish
vaku completion fish | source -

# Powershell
Write the contents of 'vaku completion powershell' and source them in your profile

# Zsh: In ~/.zshhrc
source <(vaku completion zsh)

```
vaku completion bash|fish|zsh|powershell
```

### Examples

```
vaku completion zsh
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --format string        output format: text|json (default "text")
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort                 sort output text (default true)
```

### SEE ALSO

* [vaku](vaku.md)	 - Vaku is a CLI for working with large Vault k/v secret engines


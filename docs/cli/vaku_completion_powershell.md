## vaku completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	vaku completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
vaku completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
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


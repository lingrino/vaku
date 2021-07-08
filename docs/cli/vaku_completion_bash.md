## vaku completion bash

generate the autocompletion script for bash

### Synopsis


Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:
$ source <(vaku completion bash)

To load completions for every new session, execute once:
Linux:
  $ vaku completion bash > /etc/bash_completion.d/vaku
MacOS:
  $ vaku completion bash > /usr/local/etc/bash_completion.d/vaku

You will need to start a new shell for this setup to take effect.
  

```
vaku completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --format string        output format: text|json (default "text")
  -i, --indent-char string   string used for indents (default "    ")
  -s, --sort                 sort output text (default true)
```

### SEE ALSO

* [vaku completion](vaku_completion.md)	 - generate the autocompletion script for the specified shell


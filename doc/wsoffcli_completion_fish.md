## wsoffcli completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	wsoffcli completion fish | source

To load completions for every new session, execute once:

	wsoffcli completion fish > ~/.config/fish/completions/wsoffcli.fish

You will need to start a new shell for this setup to take effect.


```
wsoffcli completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
  -n, --neo string     Neo standar by set
  -s, --serie string   serie number
```

### SEE ALSO

* [wsoffcli completion](doc/wsoffcli_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 16-Nov-2023

## wsoffcli completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	wsoffcli completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
wsoffcli completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
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

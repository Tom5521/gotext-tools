## gotext-tools completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	gotext-tools completion fish | source

To load completions for every new session, execute once:

	gotext-tools completion fish > ~/.config/fish/completions/gotext-tools.fish

You will need to start a new shell for this setup to take effect.


```
gotext-tools completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [gotext-tools completion](gotext-tools_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 14-Jul-2025

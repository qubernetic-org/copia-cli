# copia completion

## copia completion

Generate shell completion scripts

### Synopsis

Generate completion scripts for bash, zsh, fish, or powershell.

To load completions:

  # Bash
  source <(copia completion bash)

  # Zsh
  copia completion zsh > "${fpath[1]}/_copia"

  # Fish
  copia completion fish | source

  # PowerShell
  copia completion powershell | Out-String | Invoke-Expression

```
copia completion <shell> [flags]
```

### Options

```
  -h, --help   help for completion
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia](copia.md)	 - Copia CLI — source control for industrial automation


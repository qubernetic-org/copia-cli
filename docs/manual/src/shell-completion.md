# Shell Completion

Copia CLI supports tab completion for bash, zsh, fish, and PowerShell.

## Bash

```bash
# Add to ~/.bashrc
echo 'source <(copia completion bash)' >> ~/.bashrc
source ~/.bashrc
```

## Zsh

```bash
# Generate and install
copia completion zsh > "${fpath[1]}/_copia"

# Rebuild completion cache
compinit
```

Or add to `~/.zshrc`:

```bash
echo 'source <(copia completion zsh)' >> ~/.zshrc
```

## Fish

```bash
copia completion fish | source

# To persist
copia completion fish > ~/.config/fish/completions/copia.fish
```

## PowerShell

```powershell
# Add to $PROFILE
copia completion powershell | Out-String | Invoke-Expression

# To persist
copia completion powershell >> $PROFILE
```

## What Gets Completed

- Command names: `copia is<TAB>` → `copia issue`
- Subcommands: `copia issue c<TAB>` → `copia issue create` / `copia issue close` / `copia issue comment`
- Flag names: `copia issue list --s<TAB>` → `--state`
- Flag values: `copia issue list --state c<TAB>` → `closed`

# Shell Completion

Copia CLI supports tab completion for bash, zsh, fish, and PowerShell.

## Bash

```bash
# Add to ~/.bashrc
echo 'source <(copia-cli completion bash)' >> ~/.bashrc
source ~/.bashrc
```

## Zsh

```bash
# Generate and install
copia-cli completion zsh > "${fpath[1]}/_copia"

# Rebuild completion cache
compinit
```

Or add to `~/.zshrc`:

```bash
echo 'source <(copia-cli completion zsh)' >> ~/.zshrc
```

## Fish

```bash
copia-cli completion fish | source

# To persist
copia-cli completion fish > ~/.config/fish/completions/copia.fish
```

## PowerShell

```powershell
# Add to $PROFILE
copia-cli completion powershell | Out-String | Invoke-Expression

# To persist
copia-cli completion powershell >> $PROFILE
```

## What Gets Completed

- Command names: `copia-cli is<TAB>` → `copia-cli issue`
- Subcommands: `copia-cli issue c<TAB>` → `copia-cli issue create` / `copia-cli issue close` / `copia-cli issue comment`
- Flag names: `copia-cli issue list --s<TAB>` → `--state`
- Flag values: `copia-cli issue list --state c<TAB>` → `closed`

# Configuration

## Config File

Copia CLI stores credentials in `~/.config/copia/config.yml` (file permission `0600`):

```yaml
hosts:
  app.copia.io:
    token: "your-personal-access-token"
    user: "your-username"
  on-prem.company.com:
    token: "another-token"
    user: "another-user"
```

On Windows: `%USERPROFILE%\.config\copia\config.yml`

The config file is created automatically by `copia-cli auth login`.

## Authentication Precedence

Copia CLI resolves authentication in this order (highest priority first):

| Priority | Source | Use Case |
|----------|--------|----------|
| 1 | `--token` flag | One-off commands |
| 2 | `COPIA_TOKEN` env var | CI/CD pipelines |
| 3 | Config file | Daily interactive use |

## Host Resolution

The target Copia instance is resolved in this order:

| Priority | Source | Use Case |
|----------|--------|----------|
| 1 | `--host` flag | Explicit targeting |
| 2 | `COPIA_HOST` env var | CI/CD pipelines |
| 3 | Git remote URL | Auto-detection in repo |
| 4 | First config entry | Default fallback |

## Repository Context

When inside a git repository, `copia` automatically detects the owner and repo name from the git remote URL. This means you don't need to specify the repo for most commands:

```bash
cd ~/projects/my-plc-project
copia-cli issue list          # automatically uses the repo from git remote
copia-cli pr list             # same
```

Override with `--repo owner/repo` if needed.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `COPIA_TOKEN` | Authentication token (overrides config) |
| `COPIA_HOST` | Target host (overrides config) |
| `XDG_CONFIG_HOME` | Custom config directory (default: `~/.config`) |

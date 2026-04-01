<div align="center">

# Copia CLI

[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8.svg?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey.svg)]()
[![Gitea API](https://img.shields.io/badge/API-Gitea%20REST%20v1-609926.svg?logo=gitea&logoColor=white)](https://docs.gitea.com/api/)

[Copia](https://copia.io) on the command line. Built for automation engineers and CI pipelines.

[Installation](#installation) · [Usage](#usage) · [Commands](#commands) · [Configuration](#configuration) · [Roadmap](#roadmap)

</div>

---

`copia` brings the full power of [Copia](https://copia.io) — the source control platform for industrial automation — to your terminal. Modeled after [GitHub CLI (`gh`)](https://cli.github.com/), it provides a familiar interface for managing repositories, issues, pull requests, and more.

**Supported platforms:** Linux, macOS, Windows
**Supported instances:** Copia Cloud (`app.copia.io`), self-hosted Copia/Gitea instances

## Why

The official Copia Desktop app handles `clone` and `open`. That's it. There is no CLI for creating issues, opening PRs, managing releases, or querying repositories — the operations that automation engineers and CI pipelines need daily. This tool fills that gap.

## Status

**Beta.** Phase 1 (MVP) and Phase 2 (Workflow) complete. Pre-release binaries available. See the [Roadmap](#roadmap) for progress.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install qubernetic/tap/copia-cli
```

### Precompiled Binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/qubernetic/copia-cli/releases/latest).

```bash
# Linux (amd64)
curl -sL https://github.com/qubernetic/copia-cli/releases/latest/download/copia_linux_amd64.tar.gz | tar xz
sudo mv copia /usr/local/bin/

# macOS (Apple Silicon)
curl -sL https://github.com/qubernetic/copia-cli/releases/latest/download/copia_darwin_arm64.tar.gz | tar xz
sudo mv copia /usr/local/bin/
```

```powershell
# Windows (PowerShell)
Invoke-WebRequest -Uri https://github.com/qubernetic/copia-cli/releases/latest/download/copia_windows_amd64.zip -OutFile copia.zip
Expand-Archive copia.zip -DestinationPath "$env:LOCALAPPDATA\Programs\copia"
# Add to PATH manually or via System Settings
```

### Build from Source

Requires [Go 1.26+](https://go.dev/dl/).

```bash
go install github.com/qubernetic/copia-cli/cmd/copia@latest
```

## Usage

### Authenticate

```bash
# Interactive login (prompts for host and token)
copia-cli auth login

# Non-interactive (CI/agent-friendly)
copia-cli auth login --host app.copia.io --token YOUR_TOKEN

# Check auth status
copia-cli auth status
```

### Repositories

```bash
copia-cli repo list --org my-org
copia-cli repo view
copia-cli repo clone my-org/my-plc-project
copia-cli repo create my-new-repo --private
copia-cli repo delete my-org/old-repo --yes
copia-cli repo fork upstream-org/project --org my-org
```

### Issues

```bash
copia-cli issue list
copia-cli issue create --title "Fix sensor mapping" --label bug
copia-cli issue view 42
copia-cli issue close 42 --comment "Fixed in PR #7"
copia-cli issue comment 42 --body "Investigating now."
copia-cli issue edit 42 --add-label urgent --assignee john --milestone 1
```

### Pull Requests

```bash
copia-cli pr create --title "feat: add cylinder wrapper" --base develop
copia-cli pr list --state open
copia-cli pr view 7
copia-cli pr merge 7 --merge --delete-branch
copia-cli pr review 7 --approve
copia-cli pr diff 7
copia-cli pr checkout 7
```

### Releases

```bash
copia-cli release list
copia-cli release create v1.0.0 --title "Release 1.0.0" --notes "Changelog here"
copia-cli release upload v1.0.0 binary.tar.gz
copia-cli release delete v1.0.0
```

### Labels

```bash
copia-cli label list
copia-cli label create --name "critical" --color "#e11d48"
```

### JSON Output

Every list and view command supports `--json` for scripting and agent integration:

```bash
copia-cli issue list --json number,title,state
copia-cli pr view 7 --json title,mergeable,reviewers
```

## Commands

| Command | Subcommands | Description |
|---------|-------------|-------------|
| `copia-cli auth` | `login`, `logout`, `status` | Authenticate with a Copia instance |
| `copia-cli repo` | `list`, `view`, `clone`, `create`, `delete`, `fork` | Manage repositories |
| `copia-cli issue` | `list`, `create`, `view`, `close`, `comment`, `edit` | Manage issues |
| `copia-cli pr` | `list`, `create`, `view`, `merge`, `close`, `review`, `diff`, `checkout` | Manage pull requests |
| `copia-cli label` | `list`, `create` | Manage labels |
| `copia-cli release` | `list`, `create`, `delete`, `upload` | Manage releases |

> Run `copia-cli <command> --help` for detailed usage of any command.

## Configuration

### Config File

Stored at `~/.config/copia/config.yml` (file permission `0600`):

```yaml
hosts:
  app.copia.io:
    token: "your-personal-access-token"
    user: "your-username"
  on-prem.company.com:
    token: "another-token"
    user: "another-user"
```

### Authentication Precedence

| Priority | Source | Use Case |
|----------|--------|----------|
| 1 (highest) | `--token` flag | One-off commands |
| 2 | `COPIA_TOKEN` env var | CI/CD pipelines |
| 3 | Config file | Daily interactive use |

### Host Resolution

| Priority | Source | Use Case |
|----------|--------|----------|
| 1 (highest) | `--host` flag | Explicit instance targeting |
| 2 | `COPIA_HOST` env var | CI/CD pipelines |
| 3 | Git remote URL | Auto-detection in a repo directory |
| 4 | First config entry | Default fallback |

### Repository Context

When inside a git repository, `copia` automatically detects the owner and repo name from the git remote URL. Override with `--repo owner/repo`.

## Architecture

```
copia-cli/
├── cmd/copia/                # Entrypoint
├── internal/
│   ├── build/                # Version injection (ldflags)
│   ├── config/               # Config & auth management
│   └── copiacmd/             # Root command wiring
├── pkg/
│   ├── cmd/                  # Command packages (one per command group)
│   │   ├── auth/             #   login, logout, status
│   │   ├── repo/             #   list, view, clone, create, delete, fork
│   │   ├── issue/            #   list, create, view, close, comment, edit
│   │   ├── pr/               #   list, create, view, merge, close, review, diff, checkout
│   │   ├── label/            #   list, create
│   │   └── release/          #   list, create, delete, upload
│   ├── cmdutil/              # Shared CLI helpers (factory, flags, JSON)
│   ├── iostreams/            # TTY-aware I/O abstraction
│   ├── api/                  # Gitea SDK wrapper
│   └── httpmock/             # HTTP mock for testing
├── acceptance/               # CLI acceptance tests
├── docs/                     # Developer documentation
├── .goreleaser.yml           # Cross-platform release config
└── Makefile
```

## API Foundation

Copia is built on [Gitea](https://gitea.com/) and exposes a compatible REST API:

- **Base URL:** `https://{host}/api/v1/`
- **Auth:** `Authorization: token <key>` (Personal Access Token)
- **470 endpoints** across repos, issues, PRs, releases, orgs, users, and more
- **REST only** — no GraphQL
- **No anonymous access** — every call requires authentication

See [`docs/api-reference.md`](docs/api-reference.md) for the full endpoint mapping.

## Roadmap

### Phase 1 — Core (MVP)

- [x] `copia-cli auth` — login, logout, status
- [x] `copia-cli repo` — list, view, clone
- [x] `copia-cli issue` — list, create, view, close, comment
- [x] `copia-cli pr` — list, create, view, merge, close
- [x] `copia-cli label` — list, create
- [x] `--json` output on all list/view commands

### Phase 2 — Workflow

- [x] `copia-cli release` — list, create, delete, upload
- [x] `copia-cli repo` — create, delete, fork
- [x] `copia-cli pr` — review, diff, checkout
- [x] `copia-cli issue edit` — labels, assignees, milestones
- [x] Homebrew tap
- [ ] winget package (deferred to Phase 4)

### Phase 3 — Power Features

- [ ] `copia-cli api` — generic REST escape hatch
- [ ] `copia-cli search` — repos, issues
- [ ] `copia-cli org` — list, view
- [ ] `copia-cli notification` — list, read
- [ ] Tab completion (bash/zsh/powershell)

### Out of Scope

- `workflow`/`run` — Copia may not expose Gitea Actions
- `codespace`/`copilot`/`project`/`cache` — GitHub-specific, no Gitea equivalent
- GUI — this is a CLI tool

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for the development workflow, prerequisites, and project structure.

A [devcontainer](.devcontainer/) configuration is included for VS Code / Cursor — open the repo and select "Reopen in Container" to get a fully configured environment.

## Related

- **[Copia](https://copia.io)** — Source control platform for industrial automation
- **[Copia Desktop](https://copia.io/product/copia-desktop/)** — Official desktop app (clone/open only)
- **[GitHub CLI](https://cli.github.com/)** — The reference implementation this tool mirrors
- **[Gitea API Docs](https://docs.gitea.com/api/)** — Upstream API reference

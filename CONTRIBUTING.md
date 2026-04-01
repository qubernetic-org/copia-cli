# Contributing

Thank you for your interest in improving Copia CLI.

## Prerequisites

- **Docker** — required for devcontainer
- **VS Code / Cursor** — or the [devcontainer CLI](https://github.com/devcontainers/cli)
- **GitHub CLI (`gh`)** — recommended for issue/PR management

## Quick Start

The devcontainer provides a fully configured environment with Go 1.26+, gh CLI, and golangci-lint pre-installed. No local Go installation needed.

**VS Code / Cursor:**

1. Clone the repo and open it in your editor
2. "Reopen in Container" when prompted
3. All tools are pre-installed — start coding

**CLI (without VS Code):**

```bash
git clone https://github.com/qubernetic-org/copia-cli.git
cd copia-cli

# Build and start the devcontainer
npx @devcontainers/cli up --workspace-folder .

# Run commands inside the container
npx @devcontainers/cli exec --workspace-folder . make build
npx @devcontainers/cli exec --workspace-folder . make test
npx @devcontainers/cli exec --workspace-folder . ./bin/copia --version
```

> **Note:** `npx` runs the devcontainer CLI without global installation. Install globally with `npm install -g @devcontainers/cli` to use `devcontainer` directly.

## Development Workflow

This repo follows a strict Gitflow workflow. Every contribution goes through these steps:

```
1. Open a GitHub Issue describing the change
2. Fork the repo (external contributors) or create a branch directly (maintainers)
3. Create a branch from develop:
   git checkout develop
   git pull origin develop
   git checkout -b <type>/<issue>-<slug>
4. Make atomic commits using Conventional Commits format
5. Run tests:
   make test
6. Push and open a PR targeting develop:
   git push -u origin <type>/<issue>-<slug>
   gh pr create --base develop
7. Verify test plan items and check them off in the PR description
8. After merge, clean up:
   git checkout develop && git pull origin develop
   git fetch --prune && git branch -d <type>/<issue>-<slug>
```

### Branch Types

| Change type | Branch prefix | Example |
|-------------|---------------|---------|
| New feature | `feature/` | `feature/42-add-search` |
| Bug fix | `fix/` | `fix/17-broken-timeout` |
| Documentation | `docs/` | `docs/23-update-readme` |
| Hotfix | `hotfix/` | `hotfix/89-auth-crash` |

### Commit Format

```
<type>(<optional scope>): <description>
```

Use imperative mood, lowercase after colon, no period.

| Type | When to use |
|------|-------------|
| `feat` | New feature or capability |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `chore` | Maintenance, config, dependencies |
| `refactor` | Code restructure without behavior change |
| `test` | Adding or updating tests |

### Testing

Before submitting a PR:

1. **Unit tests:** `make test`
2. **Build:** `make build`
3. **Lint:** `golangci-lint run ./...` (if installed)

Integration tests require a Copia API token and run separately:

```bash
export COPIA_TEST_TOKEN=your-token
export COPIA_TEST_HOST=app.copia.io
make integration
```

## Project Structure

```
cmd/copia/          Entrypoint
internal/           Private packages (config, build, root command)
pkg/cmd/            Command packages (one per command group)
  auth/             login, logout, status
  repo/             list, view, clone, create, delete, fork
  issue/            list, create, view, close, comment, edit
  pr/               list, create, view, merge, close, review, diff, checkout
  label/            list, create
  release/          list, create, delete, upload
pkg/cmdutil/        Shared CLI helpers (factory, flags, JSON)
pkg/iostreams/      TTY-aware I/O abstraction
pkg/api/            Gitea SDK wrapper
pkg/httpmock/       HTTP mock for testing
test/integration/   Integration tests against live Copia API
```

## Guidelines

### Invariants (non-negotiable)

- **Atomic logical commits** — one commit = one change
- **Conventional Commits** format with imperative mood
- **Gitflow branch model** — `main`/`develop` protected, PR-only merges
- **Issue-driven workflow** — every branch traces to a GitHub Issue
- **TDD** — write failing tests first, then implement
- **Semantic Versioning** for releases

### Code Style

- Follow `gofmt` and `go vet` conventions
- Keep files focused — one responsibility per file
- Follow existing patterns (Options struct + NewCmd + Run function)
- Add `--json` support to every list/view command

## Reporting Bugs

Please open an issue with:

- The command you ran
- What happened (include error output)
- What you expected to happen
- Your OS and `copia --version` output

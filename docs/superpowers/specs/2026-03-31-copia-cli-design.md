# Copia CLI — Design Specification

**Date:** 2026-03-31
**Status:** Approved
**Module path:** `github.com/qubernetic/copia-cli`

---

## 1. Overview

Copia CLI is a command-line interface for [Copia](https://copia.io) — the source control platform for industrial automation. It mirrors the [GitHub CLI (`gh`)](https://cli.github.com/) in UX, command structure, and internal architecture, targeting the Gitea-compatible REST API.

**Goal:** A reliable user/agent tool for executing git operations against Copia instances — repos, issues, PRs, labels, and releases — with the same confidence and ergonomics as `gh`.

**Target users:** Automation engineers, CI/CD pipelines, AI agents.

**Supported platforms:** Linux, macOS, Windows (single cross-compiled binary).

---

## 2. Technology Stack

| Component | Choice | Rationale |
|-----------|--------|-----------|
| Language | Go 1.26+ | `gh` is Go, single binary, cross-compile, Gitea ecosystem |
| CLI framework | Cobra | Industry standard (`gh`, `kubectl`, `docker` all use it) |
| API client | Gitea Go SDK (`code.gitea.io/sdk/gitea`) | Typed structs, pagination, auth — standard Gitea endpoints fully covered |
| Config format | YAML (`gopkg.in/yaml.v3`) | `gh` convention, human-readable |
| Assertions | `testify/assert` | `gh` convention |
| Build/release | GoReleaser + Makefile | Multi-platform binary, GitHub Releases |
| Module path | `github.com/qubernetic/copia-cli` | GitHub repo path, Go convention |

---

## 3. Architecture

Follows the `gh` CLI repository structure (`github.com/cli/cli`).

```
copia-cli/
├── cmd/
│   └── copia/
│       └── main.go                     # Minimal entrypoint
├── internal/
│   ├── build/
│   │   └── build.go                    # Version/commit injection via ldflags
│   ├── config/
│   │   ├── config.go                   # YAML config management
│   │   └── auth.go                     # Token storage, host resolution
│   └── copiacmd/
│       └── root.go                     # Root command factory, global flags
├── pkg/
│   ├── cmd/                            # Domain-driven command packages
│   │   ├── auth/
│   │   │   ├── auth.go                 # Parent command
│   │   │   ├── login/login.go
│   │   │   ├── logout/logout.go
│   │   │   └── status/status.go
│   │   ├── repo/
│   │   │   ├── repo.go
│   │   │   ├── list/list.go
│   │   │   ├── view/view.go
│   │   │   └── clone/clone.go
│   │   ├── issue/
│   │   │   ├── issue.go
│   │   │   ├── list/list.go
│   │   │   ├── create/create.go
│   │   │   ├── view/view.go
│   │   │   ├── close/close.go
│   │   │   └── comment/comment.go
│   │   ├── pr/
│   │   │   ├── pr.go
│   │   │   ├── list/list.go
│   │   │   ├── create/create.go
│   │   │   ├── view/view.go
│   │   │   ├── merge/merge.go
│   │   │   └── close/close.go
│   │   └── label/
│   │       ├── label.go
│   │       ├── list/list.go
│   │       └── create/create.go
│   ├── cmdutil/
│   │   ├── factory.go                  # Shared factory: API client, config, IO
│   │   ├── flags.go                    # Common flag helpers
│   │   └── json.go                     # --json flag, field selection
│   ├── iostreams/
│   │   └── iostreams.go                # TTY-aware I/O abstraction
│   ├── api/
│   │   └── client.go                   # Gitea SDK wrapper
│   └── httpmock/
│       └── httpmock.go                 # HTTP transport mock for testing
├── acceptance/
│   ├── acceptance_test.go              # CLI blackbox tests (testscript)
│   └── testdata/                       # txtar golden files
├── script/
│   └── build.go                        # Cross-platform build helper
├── go.mod
├── go.sum
├── .goreleaser.yml
├── Makefile
└── docs/
```

### Key patterns

- **`internal/`** — Private packages: config, build info, root command wiring. Not importable by external code.
- **`pkg/`** — Public packages: commands, utilities, IO, API client. Each command group is a separate package; each subcommand is its own sub-package.
- **Factory injection** — Every command receives its dependencies (IO, API client, config) via `cmdutil.Factory`, enabling testability.
- **Domain-driven commands** — One package per subcommand (`pkg/cmd/issue/list/`), not one large file per command group.

---

## 4. Authentication

### Login flow

1. `copia auth login` — interactive prompt for host (default: `app.copia.io`) and token (masked input)
2. `copia auth login --host <HOST> --token <TOKEN>` — non-interactive (CI/agent)
3. Validation: `GET /api/v1/user` with the provided token
4. On success: write to `~/.config/copia/config.yml` with file permission `0600`
5. On failure: error message, nothing saved

### Token precedence (every API call)

| Priority | Source | Use case |
|----------|--------|----------|
| 1 (highest) | `--token` flag | One-off commands |
| 2 | `COPIA_TOKEN` env var | CI/CD pipelines |
| 3 | Config file | Daily interactive use |

### Host resolution (every API call)

| Priority | Source | Use case |
|----------|--------|----------|
| 1 (highest) | `--host` flag | Explicit instance targeting |
| 2 | `COPIA_HOST` env var | CI/CD pipelines |
| 3 | Git remote URL | Auto-detection in repo directory |
| 4 | First config entry | Default fallback |

### Config file format

```yaml
hosts:
  app.copia.io:
    token: "abc123..."
    user: "john"
  on-prem.company.com:
    token: "def456..."
    user: "jane"
```

- Path: `~/.config/copia/config.yml` (respects `XDG_CONFIG_HOME`; Windows: `%USERPROFILE%\.config\copia\config.yml`)
- File permission: `0600`, enforced on every write
- Multi-instance: supports multiple hosts simultaneously

### Security

- MVP: config file with `0600` permission (same approach as Gitea tea CLI, Docker CLI, kubectl, AWS CLI)
- Backlog: OS keyring integration (Phase 2) — `go-keyring` library for macOS Keychain, Linux libsecret, Windows Credential Manager

---

## 5. Command Design

### Command structure

```
copia <command> <subcommand> [flags]
```

### Options struct pattern (gh convention)

Every subcommand follows the same pattern:

```go
type ListOptions struct {
    IO         *iostreams.IOStreams
    Client     *gitea.Client
    Config     config.Config

    Owner      string
    Repo       string
    State      string
    Limit      int
    JSON       bool
    JSONFields []string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
    opts := &ListOptions{}
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List issues",
        RunE: func(cmd *cobra.Command, args []string) error {
            opts.IO = f.IOStreams
            opts.Client = f.Client()
            opts.Config = f.Config()
            return listRun(opts)
        },
    }
    cmd.Flags().StringVar(&opts.State, "state", "open", "Filter by state: {open|closed|all}")
    cmd.Flags().IntVar(&opts.Limit, "limit", 30, "Maximum number of items")
    cmdutil.AddJSONFlags(cmd, &opts.JSON, &opts.JSONFields)
    return cmd
}
```

### Repo context auto-detection

When inside a git repository, `owner/repo` is parsed from the git remote URL:
- `https://{host}/{owner}/{repo}.git`
- `git@{host}:{owner}/{repo}.git`

Override with `--repo owner/repo`. Error if no repo context and no `--repo` flag.

### Output

- **Human-readable** by default: tables for lists, structured text for views
- **`--json` flag** on every list/view command: JSON array/object output with field selection (whitelisted fields per command, following `gh` convention — e.g., `--json number,title,state`)
- **Exit codes:** `0` success, `1` error, `4` auth failure

---

## 6. Testing Strategy

Three layers, following the `gh` CLI test organization.

### Unit tests (colocated `*_test.go`)

- Mock IO via `iostreams` fakes
- Mock HTTP via custom `pkg/httpmock/` package (transport-level interception)
- Response fixtures in `fixtures/` directories per command
- Run with: `go test ./...`
- TDD approach: tests written before implementation

### Integration tests (colocated `*_integration_test.go`)

- Build tag: `//go:build integration`
- Real API calls against `app.copia.io` (prod instance)
- Dedicated test repo + API token via env vars (`COPIA_TEST_TOKEN`, `COPIA_TEST_HOST`)
- Every test cleans up after itself (create → assert → delete)
- Run with: `go test -tags=integration ./...`
- Skipped when env vars missing: `t.Skip("COPIA_TEST_TOKEN not set")`

### Acceptance tests (`acceptance/`)

- Build tag: `//go:build acceptance`
- CLI binary blackbox tests using `testscript` (txtar format)
- Golden file output comparison
- Run with: `go test -tags=acceptance ./acceptance`

### Makefile targets

```makefile
make test          # go test ./...
make integration   # go test -tags=integration ./...
make acceptance    # go test -tags=acceptance ./acceptance
```

---

## 7. Distribution

### MVP

- **GitHub Releases** — precompiled binaries for all platforms via GoReleaser
- **`go install`** — for Go developers
- Targets: `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`, `windows/amd64`

### Phase 2

- **Homebrew tap** (`qubernetic/tap/copia`)
- **winget** (manifest PR to `microsoft/winget-pkgs`, automated via GoReleaser)

### Versioning

- Semantic versioning: `v0.1.0`, `v1.0.0`, ...
- Git tag triggers GoReleaser in CI
- Version injection via ldflags → `internal/build/` package
- `copia --version` → `copia version 0.1.0 (commit: abc1234, built: 2026-03-31)`

---

## 8. Phase 1 MVP Scope

### Commands (18 subcommands)

| Command | Subcommands |
|---------|-------------|
| `copia auth` | `login`, `logout`, `status` |
| `copia repo` | `list`, `view`, `clone` |
| `copia issue` | `list`, `create`, `view`, `close`, `comment` |
| `copia pr` | `list`, `create`, `view`, `merge`, `close` |
| `copia label` | `list`, `create` |

### Cross-cutting concerns

- `--json` flag on all list/view commands
- Repo context auto-detection from git remote
- Auth precedence: `--token` > `COPIA_TOKEN` env > config
- Host resolution: `--host` > `COPIA_HOST` env > git remote > config
- Unified error handling + exit codes (0/1/4)

### Implementation order

1. Project skeleton (go.mod, cobra root, Makefile, CI)
2. `copia auth` (login/logout/status) — foundation for everything else
3. `copia repo` (list/view/clone) — simple read-only, validates API integration
4. `copia label` (list/create) — simple CRUD, establishes the command pattern
5. `copia issue` (list/create/view/close/comment) — full CRUD
6. `copia pr` (list/create/view/merge/close) — most complex

### Out of scope (MVP)

- `copia release` — Phase 2
- `copia repo create/delete/fork` — Phase 2
- `copia pr review/diff/checkout` — Phase 2
- `copia api` escape hatch — Phase 3
- Homebrew/winget — Phase 2
- OS keyring — Phase 2
- Tab completion — Phase 3
- Search, orgs, notifications — Phase 3
- PLC binary diff, DeviceLink — out of scope entirely

---

## 9. Development Workflow

- **Git workflow:** git-workflow skill (Gitflow, conventional commits, semantic versioning)
- **Development method:** TDD — tests first, implementation second
- **CI:** GitHub Actions — `make test` on every PR, `make integration` with secrets
- **Release:** Git tag → GoReleaser → GitHub Releases

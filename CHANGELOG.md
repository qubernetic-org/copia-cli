# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.0-rc.2] - 2026-04-02

### Added

- Release variants: .deb, .rpm packages via nfpms
- macOS universal binary (amd64+arm64 fat binary)
- linux/386, linux/arm (v6), windows/386 build targets
- `-R`/`--repo` flag for all repo-scoped commands (#105)
- Jekyll manual site matching gh CLI layout (#132)
- Command group annotations (General/Targeted) on parent commands
- Long descriptions and descriptive examples for all commands
- Platform install guides (Linux, macOS, Windows, source)
- `notification list --all` flag (#109)
- `issue list --label` flag (#111)
- `pr close --comment` and `--delete-branch` flags (#112)

### Fixed

- `auth login` failed with "file already closed" due to premature response body close (#97)
- `resp.Body.Close` before `ReadAll` in all HTTP commands (#100)
- `repo clone` failed on private repos — missing auth (#101)
- PATCH/POST commands rejected valid 201 status from Gitea API (#102)
- `--state` flag not validated in issue/pr list (#106)
- `auth status` did not error for unknown host (#107)
- `search issues` returned empty results — defaulted to open only (#108)
- `notification list` HTTP 500 without page parameter (#109)
- `--limit` accepted negative values (#110)
- GoReleaser archive format deprecation warnings (#96)

### Changed

- License changed from MIT to AGPL-3.0 + Commercial dual license (#135)
- Documentation migrated from mdBook to Jekyll (#132)
- README restructured to match gh CLI lean style (#83)
- Integration tests rewritten to exercise CLI code paths (#103)
- `docs/` reorganized to match gh CLI convention

## [0.4.0-rc.1] - 2026-04-02

### Changed

- Binary renamed from `copia` to `copia-cli` (avoids conflict with Copia Desktop)
- Organization renamed from `qubernetic-org` to `qubernetic`
- Go module path: `github.com/qubernetic/copia-cli`

### Fixed

- Auth precedence: flag > env var > config (env vars were silently ignored)
- BaseRepo detection from git remote origin (repo-scoped commands now work)
- Errors printed to stderr (were silently swallowed)
- Clone git flag injection (added `--` separator)
- Interactive login token input (bufio.Scanner + trim)
- Error messages no longer reference non-existent `--repo` flag
- Search issues uses correct per-repo endpoint
- Issue edit `--add-label` resolves label IDs by name
- JSON field names aligned to snake_case
- `splitOwnerRepo` deduplicated to `cmdutil.SplitOwnerRepo`
- `ApiOptions` renamed to `APIOptions` (Go convention)

## [0.3.0-beta.1] - 2026-04-01

### Added

- `copia api` — generic REST escape hatch with --field, --header, --method
- `copia search repos` — search repositories across the instance
- `copia search issues` — search issues with --state filter
- `copia org list` — list user's organizations
- `copia org view` — view organization details
- `copia notification list` — list unread notifications
- `copia notification read` — mark notifications as read (single or --all)
- `copia completion` — shell completion for bash, zsh, fish, powershell
- User manual website (mdBook) with auto-generated command reference
- GitHub Pages deployment workflow for manual
- `make docs` target for command reference generation

## [0.2.0-beta.1] - 2026-04-01

### Added

- `copia release list` — list releases with --json
- `copia release create` — create release with --draft, --prerelease
- `copia release delete` — delete release by tag
- `copia release upload` — upload assets to release
- `copia repo create` — create repo with --org, --private
- `copia repo delete` — delete repo with --yes confirmation
- `copia repo fork` — fork repo with optional --org
- `copia pr review` — submit review (--approve, --request-changes, --comment)
- `copia pr diff` — view PR diff output
- `copia pr checkout` — check out PR branch locally
- `copia issue edit` — edit title, body, labels, assignees, milestone
- Homebrew tap distribution (`brew install qubernetic/tap/copia`)
- Go vulnerability check (govulncheck) in CI pipeline and weekly SARIF scan
- Go mod tidy check in CI
- Go version auto-bump workflow
- Dependabot target-branch set to develop

## [0.1.0-beta.1] - 2026-04-01

### Added

- `copia auth login` — authenticate with token validation
- `copia auth logout` — remove host from config
- `copia auth status` — display hosts with token validity check
- `copia repo list` — list user/org repos with --json
- `copia repo view` — view repo details with --json
- `copia repo clone` — clone via owner/repo or URL
- `copia repo create` — create repo with --org and --private
- `copia repo delete` — delete repo with --yes confirmation
- `copia repo fork` — fork repo with optional --org
- `copia issue list` — list issues with --state, --limit, --json
- `copia issue create` — create issue with --title, --body, --label
- `copia issue view` — view issue details with --json
- `copia issue close` — close issue with optional --comment
- `copia issue comment` — add comment to issue
- `copia pr list` — list PRs with --state, --limit, --json
- `copia pr create` — create PR with --title, --body, --base, --head
- `copia pr view` — view PR details with --json
- `copia pr merge` — merge with --merge/--squash/--rebase, --delete-branch
- `copia pr close` — close a PR
- `copia label list` — list repo labels with --json
- `copia label create` — create label with --name, --color, --description
- `copia release list` — list releases with --json
- `copia release create` — create release with --draft, --prerelease
- `copia release delete` — delete release by tag
- `copia release upload` — upload assets to release
- YAML config management with multi-host support (~/.config/copia/config.yml)
- TTY-aware IOStreams abstraction for testable output
- HTTP mock registry for unit testing
- Factory dependency injection (gh CLI pattern)
- --json flag on all list/view commands
- Devcontainer with Go 1.26, gh CLI, golangci-lint, Claude Code
- GoReleaser for cross-platform releases
- GitHub Actions CI (test, lint, govulncheck, integration tests)
- CodeQL security analysis
- Dependabot for Go modules and GitHub Actions
- Go vulnerability check (govulncheck) with SARIF upload
- Go version auto-bump workflow
- Auto-close linked issues on non-default branch merges
- Integration tests against live Copia API

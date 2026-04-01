# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

# Jekyll Documentation Migration — Design Spec

**Date:** 2026-04-02
**Status:** Draft
**Closes:** #125, #126, #127, #128, #129, #130, #83, #77

## Goal

Migrate the Copia CLI manual from mdBook to Jekyll, replicating the gh CLI manual (https://cli.github.com/manual/) layout, theme, and doc generation as closely as possible. Users familiar with `gh` should find the `copia-cli` manual immediately recognizable.

## Decision Log

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Website engine | Jekyll (from mdBook) | gh CLI uses Jekyll; GitHub Pages native support |
| Hosting | Same repo (`docs/site/`) | Simpler than separate repo; project Pages sufficient |
| Doc generator | Adapt gh CLI `internal/docs/markdown.go` | 1:1 structure parity with gh manual |
| Generator cleanup | Remove GitHub-specific content, keep structure | Professional output, maintainable upstream diff |
| Command grouping | General + Targeted on parent pages | Matches gh CLI subcommand categorization |
| Guides content | Move to README, not in manual site | gh CLI manual = command reference only |
| Jekyll theme | Custom minimal (1 layout, 1 CSS, 1 config) | Full control, no gem dependencies |

## Architecture

### File Structure (After Migration)

```
copia-cli/
├── docs/
│   └── site/                          # Jekyll site root
│       ├── _config.yml                # Jekyll config
│       ├── _layouts/
│       │   └── manual.html            # Single layout: sidebar + content
│       ├── assets/
│       │   └── css/
│       │       └── style.css          # gh CLI dark theme
│       └── manual/                    # Generated pages (by gen-docs.go)
│           ├── index.md               # Landing page
│           ├── copia-cli.md           # Root command
│           ├── copia-cli_issue.md     # Parent: General/Targeted groups
│           ├── copia-cli_issue_close.md
│           └── ...                    # ~40 command pages
├── internal/
│   └── docs/
│       ├── markdown.go                # Adapted from gh CLI internal/docs/
│       └── markdown_test.go           # Tests for generator
├── script/
│   └── gen-docs.go                    # Entry point: calls internal/docs/
├── .github/workflows/
│   └── docs.yml                       # Updated: Jekyll build + deploy
├── Makefile                           # Updated: docs target
└── README.md                          # Updated: lean gh-style + guides content
```

### What Gets Deleted

```
docs/manual/                           # Entire mdBook structure
├── book.toml
├── theme/
└── src/
    ├── SUMMARY.md
    ├── index.md
    ├── installation.md
    ├── quickstart.md
    ├── configuration.md
    ├── shell-completion.md
    ├── json-output.md
    ├── ci-cd.md
    └── commands/                      # Replaced by docs/site/manual/
site/                                  # Old mdBook build output
```

### docs/ Folder Convention

Following gh CLI: `docs/` contains **internal developer documentation**, all committed.
The user-facing manual is generated into `docs/site/` which is in `.gitignore` (build artifact).

`docs/README.md` states: "This folder is used for documentation related to developing copia-cli. User docs are available at [manual link]."

Existing developer docs stay:
- `docs/api-reference.md`, `docs/authentication.md`, `docs/gh-parity.md`
- `docs/gh-cli-patterns.md`, `docs/CODEBASE_MAP.md`
- `docs/superpowers/specs/`

### Generator Flow

```
Cobra command tree
    ↓
script/gen-docs.go (entry point)
    ↓
internal/docs/markdown.go (adapted from gh CLI)
    ↓
docs/site/manual/*.md (Jekyll markdown with front matter)
    ↓
Jekyll build (CI or local)
    ↓
_site/ → GitHub Pages
```

## Component Details

### 1. Doc Generator (`internal/docs/markdown.go`)

**Source:** `github.com/cli/cli` → `internal/docs/markdown.go` (~1554 lines)

**Adaptations needed:**
- `gh` → `copia-cli` in all generated content
- Remove GitHub Enterprise references
- Remove GraphQL-specific content
- Remove `gh` GitHub-specific features (Codespaces, Copilot, etc.)
- Jekyll front matter: `layout: manual` + `permalink: /:path/:basename`
- Keep: flag rendering (`<dl>` HTML), section ordering, tree walking, `GroupedCommands()` pattern
- Keep: aliases section, JSON fields section, examples formatting

**Key functions to adapt:**
- `GenMarkdownTreeCustom()` — tree walker, minimal changes
- `genMarkdownCustom()` — per-page content, branding changes
- `printFlagsHTML()` — flag `<dl>` rendering, keep as-is

### 2. Jekyll Theme

**`_layouts/manual.html`:**
- HTML5 boilerplate
- Sidebar (nav) with command tree — generated from command structure
- Content area (main) with markdown-rendered content
- Search (optional — can add later)
- GitHub repo link in header

**`assets/css/style.css`:**
Values from Playwright audit of cli.github.com/manual/:

| Property | Value |
|----------|-------|
| Body bg | `#0d1117` |
| Content bg | `#161b22` |
| Text | `#e6edf3` |
| Links | `#2f81f7` |
| Code bg | `#010409` |
| Code block bg | `#161b22` |
| Borders | `#30363d` |
| Body font | System stack + emoji fallbacks |
| Code font | Monospace stack |
| Body font-size | 16px (content), 14px (sidebar) |
| Content padding | 48px 40px 128px |
| H1 | 32px, weight 600 |
| H2 | 24px, weight 600, margin-top 24px |
| Flag `<dt>` code bg | `#010409` |
| Flag `<dd>` color | `#8b949e` |
| Sidebar links | white (parents), `#2f81f7` (sub-commands) |
| Sidebar font | Monospace stack, 14px |

**No:** theme toggle, hamburger menu, print button, prev/next arrows.

### 3. Jekyll Config (`_config.yml`)

```yaml
title: Copia CLI Manual
description: CLI for Copia — source control for industrial automation
baseurl: /copia-cli
url: https://qubernetic.github.io
markdown: kramdown
kramdown:
  input: GFM
exclude:
  - README.md
  - Gemfile
  - Gemfile.lock
```

### 4. Sidebar Generation

The sidebar is part of `_layouts/manual.html`. Two options:

**Option A:** Static sidebar generated by `gen-docs.go` as a `_includes/sidebar.html` partial — rebuilt on each doc generation.

**Chosen: Option A** — simplest, matches gh CLI approach where sidebar is part of the build.

The sidebar structure:
```
Getting started          ← white, bold
  (link to landing page)

copia-cli                ← white, bold (root command)

api                      ← white, bold (no sub-commands)
auth                     ← white, bold
  login                  ← blue (sub-command)
  logout
  status
completion               ← white, bold
issue                    ← white, bold
  close                  ← blue
  comment
  create
  ...
```

### 5. Command Content Enrichment

Part of the generator adaptation, these features need to work:

| Feature | gh CLI | Implementation |
|---------|--------|----------------|
| Subcommand grouping | General + Targeted | `GroupedCommands()` with annotations |
| Aliases section | `### Aliases` | Detect `Aliases` field on Cobra commands |
| JSON fields section | `### JSON Fields` | Extract `validJSONFields` from commands |
| Descriptive examples | `# comment` + `$ command` | Enrich `Example` fields in Cobra commands |
| Long descriptions | Extended usage text | Add `Long` field to Cobra commands |
| Flag `<dl>` rendering | HTML definition lists | Already in gh CLI generator |

### 6. README Restructure (#83)

Migrate to gh CLI README style:
- Badges (CI, Go Report, release)
- One-liner description
- Quick install (link to platform guides)
- Basic usage (3-4 commands)
- Link to manual
- Contributing
- Related links

Guides content that was in mdBook (installation, quickstart, config, shell-completion, json-output, ci-cd) moves here or to `docs/` markdown files.

### 7. Platform Install Guides (#77)

Create in `docs/`:
- `docs/install_linux.md`
- `docs/install_macos.md`
- `docs/install_windows.md`
- `docs/install_source.md`
- `docs/releasing.md`
- `docs/README.md` (index)

### 8. CI/CD Workflow Update

**Current (`docs.yml`):**
```
Go setup → gen-docs.go → mdBook build → Pages deploy
```

**New:**
```
Go setup → gen-docs.go → Jekyll build → Pages deploy
```

Replace `jontze/action-mdbook@v4` with Jekyll build (`actions/jekyll-build-pages`).

Triggers stay the same: push to main on `docs/site/**`, `pkg/cmd/**`, workflow changes.

### 9. Local Development

```bash
# Generate command docs
make docs

# Serve locally (auto-reload)
make docs-serve
# → http://localhost:4000/copia-cli/manual/

# Or manually:
cd docs/site && bundle exec jekyll serve
```

Makefile targets:
```makefile
docs:
	go run script/gen-docs.go

docs-serve: docs
	cd docs/site && bundle exec jekyll serve --baseurl /copia-cli

docs-clean:
	rm -rf docs/site/_site docs/site/manual/copia-cli_*.md
```

**Prerequisites:** `gem install bundler jekyll` (or use GitHub Actions for CI only).

## Issue Consolidation

The Jekyll migration supersedes several existing issues:

| Issue | Action |
|-------|--------|
| #125 (mdBook theme) | **Superseded** — Jekyll theme replaces mdBook CSS |
| #126 (example comments) | **Included** — generator + Cobra Example fields |
| #127 (Long descriptions) | **Included** — Cobra Long fields + generator |
| #128 (sidebar restructure) | **Superseded** — Jekyll sidebar replaces mdBook |
| #129 (JSON fields section) | **Included** — generator handles this |
| #130 (aliases section) | **Included** — generator handles this |
| #83 (README restructure) | **Included** — part of this migration |
| #77 (install guides) | **Included** — platform docs in `docs/` |
| #75 (v1.0.0 prep) | **Partially** — badges and README part covered |

## Execution Order

### Phase 1: Generator & Jekyll Foundation
1. Copy and adapt `internal/docs/markdown.go` from gh CLI
2. Write tests for the generator
3. Update `script/gen-docs.go` to use new generator
4. Create Jekyll structure (`_config.yml`, `_layouts/manual.html`, `assets/css/style.css`)
5. Generate docs and verify locally with `jekyll serve`
6. Visual comparison with gh CLI manual using Playwright

### Phase 2: Command Content Enrichment
7. Add `GroupedCommands()` annotations to parent commands
8. Add `Long` descriptions to all commands
9. Add descriptive `Example` comments (`# comment` + `$ command`)
10. Verify aliases and JSON fields sections render correctly

### Phase 3: Docs & README Migration
11. Reorganize `docs/` folder — rename and restructure files to match gh CLI convention:
    - `docs/README.md` — index ("This folder is for development docs")
    - `docs/api-reference.md` → review/rename if needed
    - `docs/authentication.md` → review/rename if needed
    - `docs/gh-parity.md` → `docs/parity.md` or similar
    - `docs/gh-cli-patterns.md` → `docs/cli-patterns.md`
    - `docs/CODEBASE_MAP.md` → `docs/project-layout.md` (match gh CLI naming)
    - Remove or archive obsolete files
12. Create platform install guides (`docs/install_linux.md`, `docs/install_macos.md`, `docs/install_windows.md`, `docs/install_source.md`)
13. Create `docs/releasing.md`
14. Restructure README (gh CLI style, lean)

### Phase 4: CI/CD & Cleanup
15. Update `docs.yml` workflow (Jekyll build)
16. Delete mdBook structure (`docs/manual/`)
17. Add `docs/site/_site/` and `site/` to `.gitignore`
18. Update Makefile targets
19. Close consolidated issues
20. Visual regression test (Playwright side-by-side)
21. Deploy and verify on GitHub Pages

## Out of Scope

- Search functionality (can add later)
- Dark/light theme toggle (single dark theme)
- Release notes page (use GitHub Releases)
- API reference in manual (stays in `docs/api-reference.md`)

## Risks

| Risk | Mitigation |
|------|------------|
| Jekyll build failures in CI | Test locally first, pin Jekyll version |
| gh CLI generator too tightly coupled to gh | Adapt incrementally, test each change |
| Sidebar generation complexity | Start with static `_includes/sidebar.html` |
| GitHub Pages Jekyll version mismatch | Use `github-pages` gem for compatibility |

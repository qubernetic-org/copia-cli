# Jekyll Documentation Migration — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Migrate the Copia CLI manual from mdBook to Jekyll, replicating the gh CLI manual (cli.github.com/manual/) as closely as possible — same generator, layout, theme, and structure.

**Architecture:** Adapt the gh CLI `internal/docs/markdown.go` generator to produce Jekyll-compatible markdown with front matter. Create a minimal Jekyll site (`docs/site/`) with a custom dark theme matching gh CLI. The generator walks the Cobra command tree and produces one `.md` file per command with HTML flag rendering, grouped subcommands, aliases, JSON fields, and examples.

**Tech Stack:** Go (generator), Jekyll (site builder), GitHub Pages (hosting), GitHub Actions (CI/CD)

**Spec:** `docs/superpowers/specs/2026-04-02-jekyll-docs-migration-design.md`

---

## Phase 1: Generator & Jekyll Foundation

### Task 1: Clone and adapt gh CLI doc generator

**Files:**
- Create: `internal/docs/markdown.go`
- Create: `internal/docs/markdown_test.go`

- [ ] **Step 1: Clone gh CLI `internal/docs/markdown.go`**

```bash
mkdir -p internal/docs
# Fetch from gh CLI repo
curl -sL https://raw.githubusercontent.com/cli/cli/trunk/internal/docs/markdown.go -o internal/docs/markdown.go
curl -sL https://raw.githubusercontent.com/cli/cli/trunk/internal/docs/markdown_test.go -o internal/docs/markdown_test.go
```

- [ ] **Step 2: Update package and imports**

In `internal/docs/markdown.go`:
- Change `package docs` (should already be correct)
- Update module imports: remove any `github.com/cli/cli/v2/` imports
- Keep: `github.com/spf13/cobra`, `github.com/spf13/pflag`, standard library
- Remove imports for: `gh` internal packages (`pkg/cmd/root`, etc.)

- [ ] **Step 3: Remove GitHub-specific content**

Search and remove/replace in `markdown.go`:
- All references to `gh` CLI branding → `copia-cli`
- GitHub Enterprise references → remove
- GraphQL-specific content → remove
- Codespaces, Copilot, Actions references → remove
- `gh help formatting` → `copia-cli --help`
- `gh help environment` → remove or adapt
- Any URL references to `docs.github.com` → remove

- [ ] **Step 4: Adapt the `genMarkdownCustom` function**

The front matter should produce:
```yaml
---
layout: manual
permalink: /:path/:basename
---
```

Keep all section rendering logic:
- Usage/synopsis block
- Description (Long or Short)
- Subcommand groups (General/Targeted via `GroupedCommands`)
- Options (`printFlagsHTML` with `<dl>` lists)
- Inherited options
- Aliases section
- JSON fields section
- Examples section (with `{% highlight bash %}...{% endhighlight %}`)
- See also section

- [ ] **Step 5: Simplify `GroupedCommands` integration**

The gh CLI uses `root.GroupedCommands()` which depends on internal gh types. Replace with a simpler approach:

Create a helper function that checks for a Cobra command annotation:
```go
func getCommandGroups(cmd *cobra.Command) (general, targeted []*cobra.Command) {
    for _, sub := range cmd.Commands() {
        if !sub.IsAvailableCommand() || sub.IsAdditionalHelpTopicCommand() {
            continue
        }
        if sub.Annotations != nil && sub.Annotations["group"] == "targeted" {
            targeted = append(targeted, sub)
        } else {
            general = append(general, sub)
        }
    }
    return
}
```

- [ ] **Step 6: Adapt `printFlagsHTML` template**

Keep the HTML `<dl class="flags">` template as-is from gh CLI. This renders flags as:
```html
<dl class="flags">
  <dt><code>-c</code>, <code>--comment &lt;string&gt;</code></dt>
  <dd>Leave a closing comment</dd>
</dl>
```

- [ ] **Step 7: Adapt sidebar generation**

Add a function to generate `_includes/sidebar.html`:
```go
func GenSidebar(cmd *cobra.Command, dir string) error {
    // Walk command tree
    // Generate HTML nav structure
    // Write to dir/_includes/sidebar.html
}
```

The sidebar HTML should produce:
```html
<nav>
  <h5>Getting started</h5>
  <ul><li><a href="/copia-cli/manual/">Manual</a></li></ul>

  <h5>copia-cli</h5>
  <!-- For each top-level command group -->
  <h5>auth</h5>
  <ul>
    <li><a href="/copia-cli/manual/copia-cli_auth_login">login</a></li>
    <li><a href="/copia-cli/manual/copia-cli_auth_logout">logout</a></li>
    <li><a href="/copia-cli/manual/copia-cli_auth_status">status</a></li>
  </ul>
  <!-- ... -->
</nav>
```

- [ ] **Step 8: Verify the code compiles**

```bash
go build ./internal/docs/
```
Expected: no errors.

- [ ] **Step 9: Write basic tests**

In `internal/docs/markdown_test.go`, adapt the gh CLI tests:
- Test that `genMarkdownCustom` produces front matter
- Test that flags render as `<dl>` HTML
- Test that subcommands are grouped
- Test that aliases section appears for commands with aliases
- Test that `copia-cli` branding is used (no `gh` references)

```bash
go test ./internal/docs/ -v
```

- [ ] **Step 10: Commit**

```bash
git add internal/docs/
git commit -m "feat: add Jekyll doc generator adapted from gh CLI (#125)"
```

---

### Task 2: Update gen-docs.go entry point

**Files:**
- Modify: `script/gen-docs.go`

- [ ] **Step 1: Rewrite gen-docs.go to use internal/docs**

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/qubernetic/copia-cli/internal/config"
    "github.com/qubernetic/copia-cli/internal/copiacmd"
    "github.com/qubernetic/copia-cli/internal/docs"
    "github.com/qubernetic/copia-cli/pkg/cmdutil"
    "github.com/qubernetic/copia-cli/pkg/iostreams"
)

func main() {
    manualDir := "docs/site/manual"
    if len(os.Args) > 1 {
        manualDir = os.Args[1]
    }

    if err := os.MkdirAll(manualDir, 0755); err != nil {
        fmt.Fprintf(os.Stderr, "Error creating output dir: %v\n", err)
        os.Exit(1)
    }

    ios := iostreams.System()
    f := &cmdutil.Factory{
        IOStreams: ios,
        Config: func() (*config.Config, error) {
            return &config.Config{Hosts: map[string]*config.HostConfig{}}, nil
        },
    }

    rootCmd := copiacmd.NewRootCmd(f)
    rootCmd.DisableAutoGenTag = true

    linkHandler := func(name string) string {
        return "./" + strings.TrimSuffix(name, ".md")
    }

    filePrepender := func(filename string) string {
        return "---\nlayout: manual\npermalink: /:path/:basename\n---\n\n"
    }

    if err := docs.GenMarkdownTreeCustom(rootCmd, manualDir, filePrepender, linkHandler); err != nil {
        fmt.Fprintf(os.Stderr, "Error generating docs: %v\n", err)
        os.Exit(1)
    }

    // Generate sidebar
    siteDir := filepath.Dir(manualDir)
    if err := docs.GenSidebar(rootCmd, siteDir); err != nil {
        fmt.Fprintf(os.Stderr, "Error generating sidebar: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("Generated docs in %s\n", manualDir)
}
```

- [ ] **Step 2: Verify it compiles and generates**

```bash
go run script/gen-docs.go
ls docs/site/manual/
```
Expected: `copia-cli_*.md` files with Jekyll front matter.

- [ ] **Step 3: Commit**

```bash
git add script/gen-docs.go
git commit -m "refactor: update gen-docs to use internal/docs generator (#125)"
```

---

### Task 3: Create Jekyll site structure

**Files:**
- Create: `docs/site/_config.yml`
- Create: `docs/site/_layouts/manual.html`
- Create: `docs/site/assets/css/style.css`
- Create: `docs/site/Gemfile`
- Create: `docs/site/manual/index.md`

- [ ] **Step 1: Create Jekyll config**

`docs/site/_config.yml`:
```yaml
title: Copia CLI Manual
description: CLI for Copia — source control for industrial automation
baseurl: /copia-cli
url: https://qubernetic.github.io
markdown: kramdown
kramdown:
  input: GFM
exclude:
  - Gemfile
  - Gemfile.lock
  - README.md
```

- [ ] **Step 2: Create Gemfile**

`docs/site/Gemfile`:
```ruby
source "https://rubygems.org"
gem "github-pages", group: :jekyll_plugins
```

- [ ] **Step 3: Create layout template**

`docs/site/_layouts/manual.html`:
```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{ page.title | default: site.title }}</title>
  <link rel="stylesheet" href="{{ '/assets/css/style.css' | relative_url }}">
</head>
<body class="manual">
  <header>
    <nav class="d-flex flex-justify-between mx-auto px-3">
      <a href="{{ '/' | relative_url }}" class="header-logo">
        <span>CLI</span>
      </a>
      <div class="header-links">
        <a href="{{ '/manual/' | relative_url }}">Manual</a>
        <a href="https://github.com/qubernetic/copia-cli/releases">Release notes</a>
      </div>
    </nav>
  </header>

  <div class="d-flex">
    <aside class="sidebar">
      {% include sidebar.html %}
    </aside>

    <main class="main-content markdown-body">
      <div class="container-lg">
        {{ content }}
      </div>
    </main>
  </div>
</body>
</html>
```

- [ ] **Step 4: Create CSS**

`docs/site/assets/css/style.css` — the full gh CLI dark theme from our Playwright audit. Key values:

```css
/* GitHub Dark Theme for Copia CLI Manual */
* { margin: 0; padding: 0; box-sizing: border-box; }

body.manual {
  background-color: #0d1117;
  color: #e6edf3;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans",
    Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";
  font-size: 14px;
  line-height: 1.5;
}

/* Header */
header { background-color: #0d1117; border-bottom: 1px solid #30363d; padding: 16px 0; }
header nav { max-width: 1280px; display: flex; align-items: center; justify-content: space-between; }
.header-logo span { color: #ffffff; font-weight: 600; font-size: 16px; }
.header-links a { color: #8b949e; text-decoration: none; margin-left: 24px; font-size: 14px; }
.header-links a:hover { color: #e6edf3; }

/* Layout */
.d-flex { display: flex; max-width: 1280px; margin: 0 auto; }

/* Sidebar */
.sidebar {
  width: 240px; min-width: 240px; padding: 24px 16px;
  background-color: #0d1117; overflow-y: auto; max-height: calc(100vh - 60px);
  position: sticky; top: 60px;
  font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
  font-size: 14px;
}
.sidebar h5 { color: #ffffff; font-weight: 600; margin-top: 20px; margin-bottom: 4px; }
.sidebar ul { list-style: none; padding: 0; margin: 0 0 0 0; }
.sidebar li { line-height: 1.6; }
.sidebar a { color: #2f81f7; text-decoration: none; }
.sidebar a:hover { color: #58a6ff; }
.sidebar a.active { color: #ffffff; font-weight: 600; }

/* Main content */
.main-content {
  flex: 1; background-color: #161b22; padding: 48px 40px 128px;
  font-size: 16px; line-height: 1.5; min-height: calc(100vh - 60px);
}
.container-lg { max-width: 900px; }

/* Typography */
h1, h2, h3, h4, h5, h6 { color: #e6edf3; font-weight: 600; }
h1 { font-size: 32px; margin-bottom: 16px; }
h2 { font-size: 24px; margin-top: 24px; margin-bottom: 16px; border-bottom: 1px solid #30363d; padding-bottom: 0.3em; }
h3 { font-size: 20px; margin-top: 24px; margin-bottom: 16px; }
p { margin-bottom: 16px; }

/* Links */
a { color: #2f81f7; text-decoration: none; }
a:hover { color: #58a6ff; }
/* Heading anchors — white, not blue */
h1 a, h2 a, h3 a, h4 a { color: #e6edf3; }

/* Code */
code {
  background-color: #010409; border-radius: 6px; padding: 0.2em 0.4em;
  font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, monospace;
  font-size: 85%; color: #e6edf3;
}
pre { background-color: #161b22; border-radius: 6px; padding: 16px; overflow-x: auto; }
pre code { background: none; padding: 0; font-size: 13.6px; }

/* Flag definition lists */
dl.flags { margin: 0 0 16px; padding-left: 16px; }
dl.flags dt { margin-top: 16px; }
dl.flags dt code { background-color: #010409; padding: 2px 6px; border-radius: 6px; }
dl.flags dd { margin-left: 16px; color: #8b949e; margin-top: 2px; }

/* Tables */
table { border-collapse: collapse; margin-bottom: 16px; }
td, th { border: 1px solid #30363d; padding: 6px 13px; }
th { background-color: #161b22; font-weight: 600; }
tr:nth-child(2n) { background-color: #0d1117; }

/* Lists */
ul, ol { padding-left: 2em; margin-bottom: 16px; }
li { margin-bottom: 4px; }

/* Blockquotes */
blockquote { border-left: 3px solid #30363d; padding-left: 16px; color: #8b949e; }
```

- [ ] **Step 5: Create landing page**

`docs/site/manual/index.md`:
```markdown
---
layout: manual
permalink: /manual/
title: Copia CLI Manual
---

## Copia CLI manual

`copia-cli`, or `copia`, is a command-line interface for Copia for use in your terminal or your scripts.

- [Available commands](./copia-cli)
- [Usage examples](#examples)

## Installation

You can find installation instructions on our [README](https://github.com/qubernetic/copia-cli#installation).

## Configuration

- Run `copia-cli auth login` to authenticate with your Copia instance.
- Declare your aliases for often-used commands with `copia-cli alias set`.

## Examples

```bash
$ copia-cli issue list
$ copia-cli issue create --label bug
$ copia-cli repo view my-org/my-repo
```

## See also

- [copia-cli](./copia-cli)
```

- [ ] **Step 6: Create `_includes` directory for sidebar**

```bash
mkdir -p docs/site/_includes
```

The sidebar.html will be generated by gen-docs.go (Task 2, Step 1).

- [ ] **Step 7: Verify Jekyll builds locally**

```bash
cd docs/site && bundle install && bundle exec jekyll build
ls _site/manual/
```
Expected: HTML files generated in `_site/`.

- [ ] **Step 8: Verify local serve**

```bash
cd docs/site && bundle exec jekyll serve --baseurl /copia-cli
```
Open `http://localhost:4000/copia-cli/manual/` — verify dark theme, sidebar, content renders.

- [ ] **Step 9: Commit**

```bash
git add docs/site/
git commit -m "feat: add Jekyll site structure with gh CLI dark theme (#125)"
```

---

### Task 4: Visual comparison with Playwright

**Files:**
- Create: `/tmp/playwright-test-jekyll-compare.js` (temporary, not committed)

- [ ] **Step 1: Write comparison script**

Take side-by-side screenshots of:
- Landing page: gh CLI vs ours
- Parent command page (issue): gh CLI vs ours
- Leaf command page (issue close): gh CLI vs ours
- Sidebar comparison

- [ ] **Step 2: Run comparison and verify**

```bash
cd ~/.claude/skills/playwright-skill && node run.js /tmp/playwright-test-jekyll-compare.js
```

Review screenshots. Fix any CSS issues found.

- [ ] **Step 3: Commit any CSS fixes**

```bash
git add docs/site/assets/css/style.css
git commit -m "fix: adjust CSS based on visual comparison (#125)"
```

---

## Phase 2: Command Content Enrichment

### Task 5: Add command group annotations

**Files:**
- Modify: `pkg/cmd/issue/issue.go`
- Modify: `pkg/cmd/pr/pr.go`
- Modify: `pkg/cmd/repo/repo.go`
- Modify: `pkg/cmd/release/release.go`
- Modify: `pkg/cmd/org/org.go`
- Modify: `pkg/cmd/notification/notification.go`

- [ ] **Step 1: Add "group" annotations to targeted commands**

For each parent command, annotate sub-commands. Example for issue:

In `pkg/cmd/issue/issue.go`, after adding sub-commands:
```go
// General commands (create, list workflows)
// Targeted commands need annotation:
closeCmd.Annotations = map[string]string{"group": "targeted"}
commentCmd.Annotations = map[string]string{"group": "targeted"}
editCmd.Annotations = map[string]string{"group": "targeted"}
viewCmd.Annotations = map[string]string{"group": "targeted"}
```

Pattern per group:
- **issue**: General = create, list | Targeted = close, comment, edit, view
- **pr**: General = create, list | Targeted = checkout, close, diff, merge, review, view
- **repo**: General = clone, create, list | Targeted = delete, fork, view
- **release**: General = create, list | Targeted = delete, upload
- **org**: All general (list, view)
- **notification**: All general (list, read)

- [ ] **Step 2: Verify annotations compile**

```bash
go build ./...
```

- [ ] **Step 3: Regenerate docs and verify grouping**

```bash
go run script/gen-docs.go
cat docs/site/manual/copia-cli_issue.md
```
Expected: "General commands" and "Targeted commands" sections.

- [ ] **Step 4: Commit**

```bash
git add pkg/cmd/
git commit -m "feat: add command group annotations for doc generation (#127)"
```

---

### Task 6: Add Long descriptions to commands

**Files:**
- Modify: `pkg/cmd/issue/list/list.go` (and all other ~30 command files)

- [ ] **Step 1: Add Long description to issue list**

```go
cmd := &cobra.Command{
    Use:   "list",
    Short: "List issues in a repository",
    Long: `List issues in a repository.

By default, this only lists open issues. Use --state to filter.

Results are sorted by most recently updated.`,
```

- [ ] **Step 2: Repeat for all commands**

Add `Long` descriptions to every command. Focus on:
- What the command does (more detail than Short)
- Default behavior
- Notable flags
- Links to related commands

Key commands needing Long descriptions:
- `issue list`, `issue create`, `issue close`, `issue view`, `issue edit`, `issue comment`
- `pr list`, `pr create`, `pr close`, `pr view`, `pr merge`, `pr diff`, `pr review`
- `repo list`, `repo view`, `repo clone`, `repo create`, `repo delete`, `repo fork`
- `label list`, `label create`
- `release list`, `release create`, `release delete`, `release upload`
- `search repos`, `search issues`
- `org list`, `org view`
- `notification list`, `notification read`
- `api`
- `auth login`, `auth logout`, `auth status`

- [ ] **Step 3: Verify and commit**

```bash
go build ./...
go run script/gen-docs.go
git add pkg/cmd/
git commit -m "feat: add Long descriptions to all commands (#127)"
```

---

### Task 7: Enrich command examples

**Files:**
- Modify: `pkg/cmd/issue/close/close.go` (and all other ~30 command files)

- [ ] **Step 1: Update examples with comments**

For each command, update the `Example` field with `# comment` + `$ command` pattern:

```go
Example: `  # Close an issue
  copia issue close 12

  # Close with a comment
  copia issue close 12 --comment "Fixed in PR #7"`,
```

- [ ] **Step 2: Repeat for all commands**

Add descriptive comments to every command's examples.

- [ ] **Step 3: Verify and commit**

```bash
go build ./...
go run script/gen-docs.go
git add pkg/cmd/
git commit -m "feat: add descriptive comments to command examples (#126)"
```

---

## Phase 3: Docs & README Migration

### Task 8: Reorganize docs/ folder

**Files:**
- Rename: `docs/gh-parity.md` → `docs/parity.md`
- Rename: `docs/gh-cli-patterns.md` → `docs/cli-patterns.md`
- Rename: `docs/CODEBASE_MAP.md` → `docs/project-layout.md`
- Create: `docs/README.md`

- [ ] **Step 1: Rename files**

```bash
cd ~/Git/copia-cli
git mv docs/gh-parity.md docs/parity.md
git mv docs/gh-cli-patterns.md docs/cli-patterns.md
git mv docs/CODEBASE_MAP.md docs/project-layout.md
```

- [ ] **Step 2: Create docs/README.md**

```markdown
# Documentation

This folder is used for documentation related to developing `copia-cli`.

User documentation for `copia-cli` is available at
[qubernetic.github.io/copia-cli/manual/](https://qubernetic.github.io/copia-cli/manual/).

## Developer docs

- [API Reference](api-reference.md) — Gitea API endpoint mapping
- [Authentication](authentication.md) — Auth methods and config
- [CLI Patterns](cli-patterns.md) — Command implementation patterns
- [Parity](parity.md) — gh CLI feature parity tracker
- [Project Layout](project-layout.md) — Codebase architecture

## Guides

- [Installing on Linux](install_linux.md)
- [Installing on macOS](install_macos.md)
- [Installing on Windows](install_windows.md)
- [Building from Source](install_source.md)
- [Releasing](releasing.md)
```

- [ ] **Step 3: Update any internal links in renamed files**

Search the renamed files for references to old filenames and update.

- [ ] **Step 4: Commit**

```bash
git add docs/
git commit -m "chore: reorganize docs/ to match gh CLI convention (#77)"
```

---

### Task 9: Create platform install guides

**Files:**
- Create: `docs/install_linux.md`
- Create: `docs/install_macos.md`
- Create: `docs/install_windows.md`
- Create: `docs/install_source.md`
- Create: `docs/releasing.md`

- [ ] **Step 1: Create install_linux.md**

Content: deb/rpm/Homebrew/binary download instructions. Reference existing installation.md content from mdBook.

- [ ] **Step 2: Create install_macos.md**

Content: Homebrew/binary download instructions.

- [ ] **Step 3: Create install_windows.md**

Content: Binary download/scoop/future winget instructions.

- [ ] **Step 4: Create install_source.md**

Content: Go install, clone + build, cross-compilation.

- [ ] **Step 5: Create releasing.md**

Content: Tag, GoReleaser, Homebrew tap, changelog, pre-release process.

- [ ] **Step 6: Commit**

```bash
git add docs/install_*.md docs/releasing.md
git commit -m "docs: add platform install guides and release process (#77)"
```

---

### Task 10: Restructure README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Rewrite README to match gh CLI style**

Structure:
```markdown
# copia-cli

`copia-cli` is [Copia](https://copia.io)'s official command line tool. ...

## Installation

### macOS
`brew install qubernetic/tap/copia-cli`

### Linux & BSD
[See Linux install docs](./docs/install_linux.md)

### Windows
[See Windows install docs](./docs/install_windows.md)

### Build from source
[See build from source](./docs/install_source.md)

## Quick Start

```bash
copia-cli auth login
copia-cli repo list
copia-cli issue create --title "Bug" --label bug
```

## Manual

Read the [manual](https://qubernetic.github.io/copia-cli/manual/) for a comprehensive reference of all commands.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
```

- [ ] **Step 2: Add badges**

```markdown
[![CI](https://github.com/qubernetic/copia-cli/actions/workflows/ci.yml/badge.svg)](...)
[![Go Report Card](https://goreportcard.com/badge/github.com/qubernetic/copia-cli)](...)
[![Release](https://img.shields.io/github/v/release/qubernetic/copia-cli)](...)
```

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: restructure README to match gh CLI style (#83)"
```

---

## Phase 4: CI/CD & Cleanup

### Task 11: Update docs.yml workflow

**Files:**
- Modify: `.github/workflows/docs.yml`

- [ ] **Step 1: Replace mdBook with Jekyll build**

```yaml
name: Deploy Manual

on:
  push:
    branches: [main]
    paths:
      - "docs/site/**"
      - "pkg/cmd/**"
      - "internal/docs/**"
      - ".github/workflows/docs.yml"
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: true

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6

      - uses: actions/setup-go@v6
        with:
          go-version: "1.26"

      - name: Generate command docs
        run: go run script/gen-docs.go

      - name: Build Jekyll site
        uses: actions/jekyll-build-pages@v1
        with:
          source: docs/site
          destination: ./_site

      - name: Upload Pages artifact
        uses: actions/upload-pages-artifact@v4

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v5
```

- [ ] **Step 2: Commit**

```bash
git add .github/workflows/docs.yml
git commit -m "ci: update docs workflow for Jekyll build (#125)"
```

---

### Task 12: Delete mdBook structure and cleanup

**Files:**
- Delete: `docs/manual/` (entire directory)
- Modify: `.gitignore`
- Modify: `Makefile`

- [ ] **Step 1: Delete mdBook structure**

```bash
git rm -r docs/manual/
rm -rf site/
```

- [ ] **Step 2: Update .gitignore**

Add `docs/site/_site/` to `.gitignore`. The existing `site/` entry can stay for backwards compat:

```
# Build artifacts
bin/
dist/
site/
docs/site/_site/
docs/site/.jekyll-cache/
```

- [ ] **Step 3: Update Makefile**

Replace docs target and add new targets:

```makefile
docs:
	go run script/gen-docs.go

docs-serve: docs
	cd docs/site && bundle exec jekyll serve --baseurl /copia-cli

docs-clean:
	rm -rf docs/site/_site docs/site/.jekyll-cache docs/site/manual/copia-cli_*.md docs/site/_includes/sidebar.html
```

- [ ] **Step 4: Commit**

```bash
git add -A
git commit -m "chore: remove mdBook, update Makefile for Jekyll (#125)"
```

---

### Task 13: Close issues and final verification

- [ ] **Step 1: Visual regression test**

Use Playwright to take side-by-side screenshots comparing our Jekyll site with gh CLI manual. Verify:
- Landing page
- Parent command page (issue)
- Leaf command page (issue close)
- Complex command page (issue list)
- Sidebar structure
- Flag rendering
- Examples formatting

- [ ] **Step 2: Close consolidated issues**

```bash
cd ~/Git/copia-cli
gh issue close 125 --comment "Superseded by Jekyll migration"
gh issue close 126 --comment "Included in Jekyll migration"
gh issue close 127 --comment "Included in Jekyll migration"
gh issue close 128 --comment "Superseded by Jekyll migration"
gh issue close 129 --comment "Included in Jekyll migration"
gh issue close 130 --comment "Included in Jekyll migration"
gh issue close 83 --comment "Included in Jekyll migration"
gh issue close 77 --comment "Included in Jekyll migration"
```

- [ ] **Step 3: Verify GitHub Pages deploy**

After merging to main, verify https://qubernetic.github.io/copia-cli/manual/ loads with the new Jekyll site.

---

## Summary

| Phase | Tasks | Commits |
|-------|-------|---------|
| 1. Generator & Jekyll | Tasks 1-4 | 4-5 |
| 2. Content Enrichment | Tasks 5-7 | 3 |
| 3. Docs & README | Tasks 8-10 | 4 |
| 4. CI/CD & Cleanup | Tasks 11-13 | 3 |
| **Total** | **13 tasks** | **~15 commits** |

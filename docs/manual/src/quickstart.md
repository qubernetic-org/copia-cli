# Quick Start

## 1. Authenticate

Generate a Personal Access Token on your Copia instance, then:

```bash
copia-cli auth login
```

This prompts for your host and token. For CI/automation:

```bash
copia-cli auth login --host app.copia.io --token YOUR_TOKEN
```

Verify your authentication:

```bash
copia-cli auth status
```

## 2. Work with Repositories

```bash
# List your repos
copia-cli repo list

# List repos in an organization
copia-cli repo list --org my-org

# Clone a repo
copia-cli repo clone my-org/my-plc-project

# View repo details
copia-cli repo view my-org/my-plc-project
```

## 3. Manage Issues

```bash
# List open issues
copia-cli issue list

# Create an issue
copia-cli issue create --title "Fix sensor mapping" --label bug

# View issue details
copia-cli issue view 42

# Add a comment
copia-cli issue comment 42 --body "Investigating now."

# Close with a comment
copia-cli issue close 42 --comment "Fixed in PR #7"
```

## 4. Pull Requests

```bash
# Create a PR
copia-cli pr create --title "feat: add safety interlock" --base develop --head feature/safety

# List open PRs
copia-cli pr list

# Review and merge
copia-cli pr review 7 --approve
copia-cli pr merge 7 --merge --delete-branch
```

## 5. JSON Output for Scripting

Every list and view command supports `--json`:

```bash
# Get issues as JSON
copia-cli issue list --json number,title,state

# Use with jq
copia-cli repo list --json fullName,description | jq '.[].fullName'
```

## Next Steps

- [Configuration](./configuration.md) — multi-host setup, environment variables
- [Shell Completion](./shell-completion.md) — tab completion for your shell
- [CI/CD Integration](./ci-cd.md) — using copia in pipelines

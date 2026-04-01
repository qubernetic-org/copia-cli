# Quick Start

## 1. Authenticate

Generate a Personal Access Token on your Copia instance, then:

```bash
copia auth login
```

This prompts for your host and token. For CI/automation:

```bash
copia auth login --host app.copia.io --token YOUR_TOKEN
```

Verify your authentication:

```bash
copia auth status
```

## 2. Work with Repositories

```bash
# List your repos
copia repo list

# List repos in an organization
copia repo list --org my-org

# Clone a repo
copia repo clone my-org/my-plc-project

# View repo details
copia repo view my-org/my-plc-project
```

## 3. Manage Issues

```bash
# List open issues
copia issue list

# Create an issue
copia issue create --title "Fix sensor mapping" --label bug

# View issue details
copia issue view 42

# Add a comment
copia issue comment 42 --body "Investigating now."

# Close with a comment
copia issue close 42 --comment "Fixed in PR #7"
```

## 4. Pull Requests

```bash
# Create a PR
copia pr create --title "feat: add safety interlock" --base develop --head feature/safety

# List open PRs
copia pr list

# Review and merge
copia pr review 7 --approve
copia pr merge 7 --merge --delete-branch
```

## 5. JSON Output for Scripting

Every list and view command supports `--json`:

```bash
# Get issues as JSON
copia issue list --json number,title,state

# Use with jq
copia repo list --json fullName,description | jq '.[].fullName'
```

## Next Steps

- [Configuration](./configuration.md) — multi-host setup, environment variables
- [Shell Completion](./shell-completion.md) — tab completion for your shell
- [CI/CD Integration](./ci-cd.md) — using copia in pipelines

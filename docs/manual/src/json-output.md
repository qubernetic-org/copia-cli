# JSON Output

Every list and view command supports `--json` for machine-readable output. This makes `copia` ideal for scripting, CI/CD pipelines, and AI agent integration.

## Basic Usage

```bash
# List issues as JSON with selected fields
copia-cli issue list --json number,title,state

# View repo details as JSON
copia-cli repo view my-org/my-repo --json fullName,description,private
```

## Available Fields

Fields vary by command. Use `copia-cli <command> <subcommand> --help` to see available JSON fields.

### Common Fields by Command

| Command | Fields |
|---------|--------|
| `repo list` | fullName, description, private, updatedAt |
| `repo view` | fullName, description, private, defaultBranch, stars, forks, openIssues |
| `issue list` | number, title, state, labels, updatedAt |
| `issue view` | number, title, body, state, author, labels, createdAt, comments |
| `pr list` | number, title, state, author, base, head, updatedAt |
| `pr view` | number, title, body, state, mergeable, author, base, head, createdAt |
| `label list` | name, color, description |
| `release list` | tagName, name, draft, prerelease, publishedAt |
| `org list` | username, fullName, description |
| `search repos` | fullName, description, htmlUrl |

## Piping with jq

```bash
# Get just repo names
copia-cli repo list --json fullName | jq -r '.[].full_name'

# Count open issues
copia-cli issue list --json number | jq length

# Filter PRs by author
copia-cli pr list --json number,title,author | jq '.[] | select(.user.login == "john")'
```

## Using in Scripts

```bash
#!/bin/bash
# Close all issues with "wontfix" label
for num in $(copia-cli issue list --json number,labels | jq -r '.[] | select(.labels[].name == "wontfix") | .number'); do
  copia-cli issue close "$num" --comment "Closing as wontfix"
done
```

## Using with the API Command

The `copia-cli api` command always returns JSON:

```bash
# Raw API call with pretty-printed JSON
copia-cli api /user

# Create issue via API
copia-cli api -X POST /repos/my-org/my-repo/issues --field title="Bug report"
```

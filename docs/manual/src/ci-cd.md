# CI/CD Integration

Copia CLI is designed for use in CI/CD pipelines. Authentication via environment variables, `--json` output for parsing, and non-interactive mode make it pipeline-friendly.

## Authentication in CI

Use environment variables — never store tokens in code:

```yaml
# GitHub Actions
env:
  COPIA_TOKEN: ${{ secrets.COPIA_TOKEN }}
  COPIA_HOST: app.copia.io
```

```bash
# Any CI system
export COPIA_TOKEN="your-token"
export COPIA_HOST="app.copia.io"
```

Or pass directly:

```bash
copia --token "$COPIA_TOKEN" --host "$COPIA_HOST" issue list
```

## GitHub Actions Example

```yaml
name: Copia Integration
on: push

jobs:
  copia:
    runs-on: ubuntu-latest
    steps:
      - name: Install Copia CLI
        run: |
          curl -sL https://github.com/qubernetic-org/copia-cli/releases/latest/download/copia_linux_amd64.tar.gz | tar xz
          sudo mv copia /usr/local/bin/

      - name: List open issues
        env:
          COPIA_TOKEN: ${{ secrets.COPIA_TOKEN }}
          COPIA_HOST: app.copia.io
        run: copia issue list --json number,title,state

      - name: Create release
        env:
          COPIA_TOKEN: ${{ secrets.COPIA_TOKEN }}
          COPIA_HOST: app.copia.io
        run: |
          copia release create "v${{ github.ref_name }}" \
            --title "Release ${{ github.ref_name }}" \
            --notes "Automated release"
```

## Common CI Patterns

### Comment on Issue After Deploy

```bash
copia issue comment 42 --body "Deployed to staging at $(date)"
```

### Create Issue on Failure

```bash
copia issue create \
  --title "CI failure: $CI_JOB_NAME" \
  --body "Pipeline failed at $(date). See logs: $CI_JOB_URL" \
  --label bug
```

### Close Issues Referenced in Commits

```bash
# Parse commit messages for "Fixes #N" and close them
git log --oneline HEAD~5..HEAD | grep -oP 'Fixes #\K\d+' | while read num; do
  copia issue close "$num" --comment "Closed by CI deploy"
done
```

### Check PR Status Before Deploy

```bash
# Get PR mergeable status
MERGEABLE=$(copia pr view 7 --json mergeable | jq -r '.mergeable')
if [ "$MERGEABLE" != "true" ]; then
  echo "PR is not mergeable"
  exit 1
fi
```

## Tips

- Always use `--json` in scripts — text output may change between versions
- Set `COPIA_TOKEN` and `COPIA_HOST` as repository/project secrets
- Use `copia api` for endpoints not covered by dedicated commands
- All commands exit with code `0` on success, `1` on error

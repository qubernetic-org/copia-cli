# Copia CLI Manual

`copia` is a command-line interface for [Copia](https://copia.io) — the source control platform for industrial automation. It provides a familiar, `gh`-like interface for managing repositories, issues, pull requests, releases, and more from the terminal.

## Who Is This For?

- **Automation engineers** managing PLC, HMI, and SCADA projects on Copia
- **DevOps/CI pipelines** that need to automate Copia operations
- **AI agents** that interact with Copia via structured JSON output

## Features

- **11 command groups**: auth, repo, issue, pr, label, release, search, org, notification, api, completion
- **35+ subcommands** covering the full Copia workflow
- **`--json` output** on every list and view command for scripting
- **Multi-instance support** for Copia Cloud and self-hosted instances
- **Cross-platform** binaries for Linux, macOS, and Windows

## Quick Example

```bash
# Authenticate
copia-cli auth login --host app.copia.io --token YOUR_TOKEN

# List your repos
copia-cli repo list

# Create an issue
copia-cli issue create --title "Fix sensor mapping" --label bug

# Open a PR
copia-cli pr create --title "feat: add safety interlock" --base develop

# Merge it
copia-cli pr merge 7 --merge --delete-branch
```

## Getting Help

- Run `copia-cli --help` for a list of commands
- Run `copia-cli <command> --help` for command-specific help
- [GitHub Issues](https://github.com/qubernetic/copia-cli/issues) for bug reports and feature requests

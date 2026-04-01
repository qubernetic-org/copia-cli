# copia-cli repo view

## copia-cli repo view

View a repository

```
copia-cli repo view [<owner/repo>] [flags]
```

### Examples

```
  copia repo view
  copia repo view my-org/my-repo
  copia repo view --json fullName,description
```

### Options

```
  -h, --help           help for view
      --json strings   Output JSON with selected fields: [full_name description private default_branch stars forks open_issues_count]
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia-cli repo](copia-cli_repo.md)	 - Manage repositories


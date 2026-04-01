# copia repo view

## copia repo view

View a repository

```
copia repo view [<owner/repo>] [flags]
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
      --json strings   Output JSON with selected fields: [fullName description private defaultBranch stars forks openIssues]
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia repo](copia_repo.md)	 - Manage repositories


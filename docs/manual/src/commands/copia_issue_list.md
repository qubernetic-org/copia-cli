# copia issue list

## copia issue list

List issues in a repository

```
copia issue list [flags]
```

### Examples

```
  copia issue list
  copia issue list --state closed
  copia issue list --json number,title,state
```

### Options

```
  -h, --help           help for list
      --json strings   Output JSON with selected fields: [number title state labels updatedAt]
  -L, --limit int      Maximum number of issues (default 30)
  -s, --state string   Filter by state: {open|closed|all} (default "open")
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia issue](copia_issue.md)	 - Manage issues


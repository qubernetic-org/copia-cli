# copia pr list

## copia pr list

List pull requests

```
copia pr list [flags]
```

### Examples

```
  copia pr list
  copia pr list --state closed
  copia pr list --json number,title,state
```

### Options

```
  -h, --help           help for list
      --json strings   Output JSON with selected fields: [number title state author base head updatedAt]
  -L, --limit int      Maximum number of pull requests (default 30)
  -s, --state string   Filter by state: {open|closed|all} (default "open")
```

### Options inherited from parent commands

```
      --host string    Target Copia host
      --token string   Authentication token
```

### SEE ALSO

* [copia pr](copia_pr.md)	 - Manage pull requests

